package moodle

import (
	"errors"
	"io"
	"mime/multipart"
	"moodle/config"
	"moodle/internal/core/domain"
	"moodle/internal/core/port"
	"moodle/pkg/apperrors"
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

// validate if Moodle.WebSiteName is unique
func (s *Service) ValidateSitename(sitename string) bool {
	moodles, err := s.mongoRepository.GetSiteName(sitename)
	if err != nil {
		return false
	}
	if len(moodles) > 0 {
		return false
	}
	return true
}

func (s *Service) Create(moodle domain.Moodle) (domain.Moodle, *apperrors.AppError) {
	// check if sitename is exist then return error
	if !s.ValidateSitename(moodle.WebSiteName) {
		return moodle, apperrors.Conflict("sitename is already exist", errors.New("sitename is already exist"))
	}
	var moodle_ domain.Moodle
	moodle_.Id = domain.NewMoodleID()
	moodle_.Email = moodle.Email
	moodle_.Name = moodle.WebSiteName
	moodle_.LbName = "kube_service" + "_" + config.CLUSTERID + "_" + moodle_.Id + "_moodle-service"
	moodle_.Ip = "Provisioning"
	moodle_.WebSiteName = moodle.WebSiteName + ".lms.bizflycloud.vn"
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
	// Insert moodle to mongodb
	err = s.mongoRepository.Create(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create moodle error", err)
	}
	// lb tracking
	err = s.mongoRepository.CreateLbTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create lb tracking error", err)
	}
	// maria tracking
	err = s.mongoRepository.CreateMariaTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create maria tracking error", err)
	}
	// mail tracking
	err = s.mongoRepository.CreateMailTracking(moodle_)
	if err != nil {
		return moodle, apperrors.Internal("create mail tracking error", err)
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
	err := s.mongoRepository.Delete(moodleID)
	if err != nil {
		return err
	}
	return nil
}

// Update logo
// Save logo into /tmp/moodle/{moodleID}/logo/{filename}
// Track logo into mongodb
func (s *Service) UpdateLogo(moodleID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError) {
	err := SaveFiles("/tmp/moodle/"+moodleID+"/logo/"+header.Filename, file)
	if err != nil {
		return nil, apperrors.Internal("save file error", err)
	}
	// update UpdateLogo
	var logotracking domain.LogoTracking
	logotracking.MoodleId = moodleID
	logotracking.FilePath = "/tmp/moodle/" + moodleID + "/logo/" + header.Filename
	logotracking.Status = "pending"
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
	// update UpdateFavicon
	var favicontracking domain.FaviconTracking
	favicontracking.MoodleId = moodleID
	favicontracking.FilePath = "/tmp/moodle/" + moodleID + "/favicon/" + header.Filename
	favicontracking.Status = "Pending"
	favicontracking.CreatedAt = time.Now()
	favicontracking.UpdatedAt = time.Now()
	err = s.mongoRepository.UpdateFavicon(favicontracking)
	if err != nil {
		return nil, apperrors.Internal("update favicon error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// update vision image
// save image into /tmp/moodle/{moodleID}/vision/{filename}
// track image into mongodb
func (s *Service) UpdateVisionImage(moodleID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError) {
	err := SaveFiles("/tmp/moodle/"+moodleID+"/vision/"+header.Filename, file)
	if err != nil {
		return nil, apperrors.Internal("save file error", err)
	}
	// update UpdateVisionImage
	var visiontracking domain.VisionImageTracking
	visiontracking.MoodleId = moodleID
	visiontracking.FilePath = "/tmp/moodle/" + moodleID + "/vision/" + header.Filename
	visiontracking.Status = "Pending"
	visiontracking.CreatedAt = time.Now()
	visiontracking.UpdatedAt = time.Now()
	err = s.mongoRepository.UpdateVisionImage(visiontracking)
	if err != nil {
		return nil, apperrors.Internal("update vision image error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// save multiple files from source to destination
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

// update banner image 
// save image into /tmp/moodle/{moodleID}/banner/{filename}
// track image into mongodb
func (s *Service) UpdateBannerImage(moodleID, bannerID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError) {
	err := SaveFiles("/tmp/moodle/"+moodleID+"/banner/"+bannerID+"/"+header.Filename, file)
	if err != nil {
		return nil, apperrors.Internal("save file error", err)
	}
	// update UpdateBannerImage
	var bannertracking domain.BannerImageTracking
	bannertracking.MoodleId = moodleID
	bannertracking.BannerId = bannerID
	bannertracking.FilePath = "/tmp/moodle/" + moodleID + "/banner/" + bannerID + "/" + header.Filename
	bannertracking.Status = "Pending"
	bannertracking.CreatedAt = time.Now()
	bannertracking.UpdatedAt = time.Now()
	err = s.mongoRepository.UpdateBannerImage(bannertracking)
	if err != nil {
		return nil, apperrors.Internal("update banner image error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// update video url
// track video url into mongodb
func (s *Service) UpdateVideoURL(moodleID, url string) (map[string]string, *apperrors.AppError) {
	// update UpdateVideoUrl
	var videotracking domain.VideoURLTracking
	videotracking.MoodleId = moodleID
	videotracking.VideoURL = url
	videotracking.Status = "Pending"
	videotracking.CreatedAt = time.Now()
	videotracking.UpdatedAt = time.Now()
	err := s.mongoRepository.UpdateVideoURL(videotracking)
	if err != nil {
		return nil, apperrors.Internal("update video url error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// update vision content
// track vision content into mongodb
func (s *Service) UpdateVisionContent(content domain.VisionContentTracking) (map[string]string, *apperrors.AppError) {
	// update UpdateVisionContent
	content.Status = "Pending"
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()
	err := s.mongoRepository.UpdateVisionContent(content)
	if err != nil {
		return nil, apperrors.Internal("update vision content error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// update banner slogan
// track banner slogan into mongodb
func (s *Service) UpdateBannerSlogan(slogan domain.BannerSloganTracking) (map[string]string, *apperrors.AppError) {
	// update UpdateBannerSlogan
	slogan.Status = "Pending"
	slogan.CreatedAt = time.Now()
	slogan.UpdatedAt = time.Now()
	err := s.mongoRepository.UpdateBannerSlogan(slogan)
	if err != nil {
		return nil, apperrors.Internal("update banner slogan error", err)
	}
	return map[string]string{"message": "success"}, nil
}

// update pre_installed_course
// track pre_installed_course into mongodb
func (s *Service) UpdatePreInstalledCourse(moodle domain.Moodle) (map[string]string, *apperrors.AppError) {
	// update UpdatePreInstalledCourse
	moodle.Status = "Pending"
	moodle.CreatedAt = time.Now()
	moodle.UpdatedAt = time.Now()
	err := s.mongoRepository.UpdatePreInstalledCourse(moodle)
	if err != nil {
		return nil, apperrors.Internal("update pre_installed_course error", err)
	}
	return map[string]string{"message": "success"}, nil
}
