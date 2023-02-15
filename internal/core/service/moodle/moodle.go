package moodle

import (
	"errors"
	"moodle/config"
	"moodle/internal/core/domain"
	"moodle/internal/core/port"
	"moodle/pkg/apperrors"
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

