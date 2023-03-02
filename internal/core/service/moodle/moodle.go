package moodle

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"moodle/config"
	"moodle/internal/core/domain"
	"moodle/internal/core/port"
	"moodle/pkg/apperrors"
	"moodle/pkg/helpers"
	"os"
	"path/filepath"
	"strings"
	"time"
)


type Service struct {
	mongoRepository port.MongoRepository
	mariaRepository port.MariaRepository
	k8sRepository port.K8sRepository
}

func NewService(mongoRepository port.MongoRepository, mariaRepository port.MariaRepository, k8sRepository port.K8sRepository) *Service {
	return &Service{
		mongoRepository: mongoRepository,
		mariaRepository: mariaRepository,
		k8sRepository: k8sRepository,
	}
}

// save file
func SaveFiles(dst string, file multipart.File) error {
	// create folder
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}
	// create file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	// copy file
	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}
	return nil
}

// validate if Moodle.WebSiteName is unique
func (s *Service) ValidateSitename(websitename string) bool {
	moodles, err := s.mongoRepository.GetSiteName(websitename)
	if err != nil {
		return false
	}
	if len(moodles) > 0 {
		return false
	}
	return true
}

func (s *Service) Create(moodle domain.Moodle) (domain.Moodle, *apperrors.AppError) {
	// check if email have one moodle then block to create new moodle
	moodles, err := s.mongoRepository.GetAll(moodle.Email)
	if err != nil {
		return moodle, apperrors.Internal("get all moodle by email error", err)
	}
	if len(moodles) > 0 {
		return moodle, apperrors.Conflict("your account allowed to create 1 moodle site only", errors.New("your account allowed to create 1 moodle site only"))
	}
	var moodle_ domain.Moodle
	moodle_.Id = domain.NewMoodleID()
	moodle_.Email = moodle.Email
	moodle_.Name = moodle.WebSiteName
	moodle_.LbName = "kube_service" + "_" + config.CLUSTERID + "_" + moodle_.Id + "_moodle-service"
	moodle_.Ip = "14.225.36.146"
	moodle_.WebSiteName = moodle.WebSiteName + ".lms.bizflycloud.vn"
	// check if websitename is exist then return error
	if !s.ValidateSitename(moodle_.WebSiteName) {
		return moodle, apperrors.Conflict("websitename is already exist", errors.New("websitename is already exist"))
	}
	moodle_.PreInstalledCourse = moodle.PreInstalledCourse
	moodle_.PackageName = moodle.PackageName
	moodle_.AutoScale = moodle.AutoScale
	moodle_.DocumentsStorageExtra = moodle.DocumentsStorageExtra
	moodle_.Status = "Creating"
	moodle_.CreatedAt = time.Now()
	moodle_.UpdatedAt = time.Now()
	// get package by package name
	package_, err := s.mongoRepository.GetPackage(moodle_.PackageName)
	if err != nil {
		return moodle, apperrors.InvalidInput("package name is not exist", err)
	}
	moodle_.Packages = package_
	// create moodle mariadb
	err = s.mariaRepository.CreateDB(strings.ReplaceAll(moodle_.Id, "-", "_"))
	if err != nil {
		return moodle, apperrors.Internal("create mariadb error", err)
	}
	// create namespace
	err = s.k8sRepository.CreateNamespace(moodle_.Id)
	if err != nil {
		return moodle, apperrors.Internal("create namespace error", err)
	}
	// Apply pvc
	err = s.k8sRepository.ApplyPVC(moodle_.Id)
	if err != nil {
		return moodle, apperrors.Internal("apply pvc error", err)
	}
	// Apply statefulset
	err = s.k8sRepository.ApplyStatefulSet(moodle_.Id)
	if err != nil {
		return moodle, apperrors.Internal("apply statefulset error", err)
	}
	// Apply service
	err = s.k8sRepository.ApplyService(moodle_.Id)
	if err != nil {
		return moodle, apperrors.Internal("apply service error", err)
	}
	// apply ingress
	err = s.k8sRepository.ApplyIngress(moodle_.Id, moodle_.WebSiteName)
	if err != nil {
		return moodle, apperrors.Internal("apply ingress error", err)
	}
	// Insert moodle to mongodb
	err = s.mongoRepository.Create(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create moodle error", err)
	}
	// maria tracking
	err = s.mongoRepository.CreateMariaTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create maria tracking error", err)
	}
	// create sitename tracking
	err = s.mongoRepository.CreateSiteNameTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create sitename tracking error", err)
	}
	// mail tracking
	err = s.mongoRepository.CreateMailTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create mail tracking error", err)
	}
	// create moodle tracking
	err = s.mongoRepository.CreateMoodleTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create moodle tracking error", err)
	}
	// create logo tracking
	defaultLogoUrl := fmt.Sprintf("http://%s/pluginfile.php/1/theme_edumy/headerlogo1/1676573021/Logo mới Bizfly Cloud-01.png", moodle_.WebSiteName)
	err = s.mongoRepository.CreateLogoTracking(moodle_, defaultLogoUrl)
	if err != nil {
		return moodle, apperrors.Internal("create logo tracking error", err)
	}
	// create favicon tracking
	defaultFaviconUrl := fmt.Sprintf("http://%s/pluginfile.php/1/theme_edumy/favicon/1676573021/z3665638475480_d6dab64f97b26f1c73cd539411ae9990.jpg", moodle_.WebSiteName)
	err = s.mongoRepository.CreateFaviconTracking(moodle_, defaultFaviconUrl)
	if err != nil {
		return moodle, apperrors.Internal("create favicon tracking error", err)
	}
	// preinstalled course tracking
	err = s.mongoRepository.CreatePreInstalledCourseTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create preinstalled course tracking error", err)
	}
	// moodle config
	err = s.mongoRepository.CreateMoodleConfig(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create moodle config error", err)
	}
	// delete row from 

	return moodle_, nil
}

// List moodles
func (s *Service) List(email string) ([]domain.Moodle, error) {
	moodles, err := s.mongoRepository.GetAll(email)
	if err != nil {
		return nil, err
	}
	return moodles, nil
}

// get moodle by moodleID from moodles collection in mongodb
func (s *Service) Get(moodleID string) (domain.Moodle, error) {
	moodle, err := s.mongoRepository.Get(moodleID)
	if err != nil {
		return domain.Moodle{}, err
	}
	return moodle, nil
}

// search moodle by name from moodles collection in mongodb
func (s *Service) Search(email, name string) ([]domain.Moodle, error) {
	moodles, err := s.mongoRepository.Search(email, name)
	if err != nil {
		return nil, err
	}
	return moodles, nil
}


// List courses
func (s *Service) ListCourses() ([]domain.Course, error) {
	courses, err := s.mongoRepository.ListCourses()
	if err != nil {
		return nil, err
	}
	return courses, nil
}

// list packages
func (s *Service) ListPackages() ([]domain.Package, error) {
	packages, err := s.mongoRepository.ListPackages()
	if err != nil {
		return nil, err
	}
	return packages, nil
}

// Delete moodle
func (s *Service) Delete(moodleID string) error {
	// delete namespace
	err := s.k8sRepository.DeleteNamespace(moodleID)
	if err != nil {
		return err
	}
	// mail tracking
	moodle, err := s.mongoRepository.Get(moodleID)
	var mailTracking domain.MailTracking
	mailTracking.MoodleId = moodle.Id
	mailTracking.Email = moodle.Email
	mailTracking.IsSent = false
	mailTracking.Type = "Delete"
	mailTracking.CreatedAt = time.Now()
	mailTracking.UpdatedAt = time.Now()
	err = s.mongoRepository.DeleteMailTracking(mailTracking)
	if err != nil {
		return err
	}
	// update status deleting on moodles collection
	err = s.mongoRepository.UpdateMoodleStatus(moodleID, "Deleting")
	if err != nil {
		return err
	}

	// get current maria tracking
	mariaTracking, err := s.mongoRepository.GetMariaTracking(moodleID)
	// update status in maria_tracing to "deleting"
	mariaTracking.DbStatus = "Deleting"
	err = s.mongoRepository.UpdateMariaTracking(mariaTracking)
	if err != nil {
		return err
	}
	return nil
}

// Update logo
// Save logo into /tmp/moodle/{moodleID}/logo/{filename}
// Track logo into mongodb
func (s *Service) UpdateLogo(moodleID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError) {
	// get sitename
	moodle, err := s.mongoRepository.Get(moodleID)
	if err != nil {
		return nil, apperrors.Internal("get moodle error", err)
	}
	err = SaveFiles("/tmp/moodle/"+moodleID+"/logo/"+header.Filename, file)
	if err != nil {
		return nil, apperrors.Internal("save file error", err)
	}
	// update UpdateLogo
	var logotracking domain.LogoTracking
	logotracking.MoodleId = moodleID
	logotracking.FilePath = "/tmp/moodle/" + moodleID + "/logo/" + header.Filename
	// resize logo
	_, err = helpers.ResizePng(logotracking.FilePath, 200, 200)
	if err != nil {
		return nil, apperrors.Internal("resize favicon error", err)
	}
	logotracking.Status = "Pending"
	logotracking.URL = fmt.Sprintf("http://%s/pluginfile.php/1/theme_edumy/headerlogo1/1676573021/Logo mới Bizfly Cloud-01.png", moodle.WebSiteName)
	logotracking.UpdatedAt = time.Now()
	logotracking.CreatedAt = time.Now()
	err = s.mongoRepository.UpdateLogo(logotracking)
	if err != nil {
		return nil, apperrors.Internal("update logo error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// Update favicon
// Save logo into /tmp/moodle/{moodleID}/favicon/{filename}
// Track logo into mongodb
func (s *Service) UpdateFavicon(moodleID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError) {
	err := SaveFiles("/tmp/moodle/"+moodleID+"/favicon/"+header.Filename, file)
	if err != nil {
		return nil, apperrors.Internal("save file error", err)
	}


	// get sitename
	moodle, err := s.mongoRepository.Get(moodleID)
	if err != nil {
		return nil, apperrors.Internal("get moodle error", err)
	}
	if err != nil {
		return nil, apperrors.Internal("save file error", err)
	}
	// update UpdateFavicon
	var favicontracking domain.FaviconTracking
	favicontracking.MoodleId = moodleID
	favicontracking.FilePath = "/tmp/moodle/" + moodleID + "/favicon/" + header.Filename
	// check if file is ico do not resize
	// resize favicon
	if header.Filename[len(header.Filename)-3:] != "ico" {
			_, err = helpers.ResizeImage(favicontracking.FilePath, 32, 32)
			if err != nil {
				return nil, apperrors.Internal("resize favicon error", err)
		}
	}
	if header.Filename[len(header.Filename)-3:] == "ico" {
			_, err = helpers.ResizeIco(favicontracking.FilePath, 32, 32)
			if err != nil {
				return nil, apperrors.Internal("resize ico favicon error", err)
		}
	}
	favicontracking.Status = "Pending"
	favicontracking.URL = fmt.Sprintf("http://%s/pluginfile.php/1/theme_edumy/favicon/1676573021/z3665638475480_d6dab64f97b26f1c73cd539411ae9990.jpg", moodle.WebSiteName)
	favicontracking.CreatedAt = time.Now()
	favicontracking.UpdatedAt = time.Now()
	err = s.mongoRepository.UpdateFavicon(favicontracking)
	if err != nil {
		return nil, apperrors.Internal("update favicon error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// update sitename
func (s *Service) UpdateSiteName(sitenametracking domain.SitenameTracking) (map[string]string, *apperrors.AppError) {
	// check if sitename is exist
	if !s.ValidateSitename(sitenametracking.SiteName) {
		return nil, apperrors.Conflict("sitename is already exist", errors.New("sitename is already exist"))
	}
	// udpate sitename tracking
	var sitenametracking_ domain.SitenameTracking
	sitenametracking_.IsUpdate = false
	sitenametracking_.SiteName = sitenametracking.SiteName
	sitenametracking_.MoodleId = sitenametracking.MoodleId
	sitenametracking_.UpdatedAt = time.Now()
	sitenametracking_.DbName = strings.ReplaceAll(sitenametracking.MoodleId, "-", "_")
	err := s.mongoRepository.UpdateSitenameTracking(sitenametracking_)
	if err != nil {
		return nil, apperrors.Internal("update sitename tracking error", err)
	}
	// update ingress
	// err = s.k8sRepository.UpdateIngress(sitenametracking.MoodleId, sitenametracking.SiteName+".lms.bizflycloud.vn")
	// if err != nil {
	// 	return nil, apperrors.Internal("update ingress error", err)
	// }
	return map[string]string{"message": "success"}, nil
}

// update pre_installed_course
// track pre_installed_course into mongodb
func (s *Service) UpdatePreInstalledCourse(moodle domain.Moodle) (map[string]string, *apperrors.AppError) {
	// get current moodle
	moodle_, err := s.mongoRepository.Get(moodle.Id)
	if err != nil {
		return nil, apperrors.Internal("get moodle error", err)
	}
	preInstalledCourse_ := moodle_.PreInstalledCourse

	// update UpdatePreInstalledCourse at last
	moodle.PreInstalledCourse = append(preInstalledCourse_, moodle.PreInstalledCourse...)
	// set a slice
	moodle.PreInstalledCourse = helpers.SetInt(moodle.PreInstalledCourse)

	moodle.CreatedAt = time.Now()
	moodle.UpdatedAt = time.Now()
	err = s.mongoRepository.UpdatePreInstalledCourse(moodle)
	if err != nil {
		return nil, apperrors.Internal("update pre_installed_course error", err)
	}
	// TODO update course to mariadb
	// preInstalledCourse tracking
	var preInstalledCourseTracking domain.PreInstalledCourseTracking
	preInstalledCourseTracking.MoodleId = moodle.Id
	preInstalledCourseTracking.CourseId = moodle.PreInstalledCourse
	preInstalledCourseTracking.Status = "Pending"
	preInstalledCourseTracking.UpdatedAt = time.Now()
	preInstalledCourseTracking.CreatedAt = time.Now()
	err = s.mongoRepository.UpdatePreInstalledCourseTracking(preInstalledCourseTracking)
	if err != nil {
		return nil, apperrors.Internal("update pre_installed_course tracking error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// get logo url from logo tracking
func (s *Service) GetLogo(moodleID string) (map[string]string, *apperrors.AppError) {
	logo, err := s.mongoRepository.GetLogoTracking(moodleID)
	if err != nil {
		return nil, apperrors.Internal("get logo error", err)
	}
	return map[string]string{"url": logo.URL}, nil
}

// get favicon url from favicon tracking
func (s *Service) GetFavicon(moodleID string) (map[string]string, *apperrors.AppError) {
	favicon, err := s.mongoRepository.GetFaviconTracking(moodleID)
	if err != nil {
		return nil, apperrors.Internal("get favicon error", err)
	}
	return map[string]string{"url": favicon.URL}, nil
}
