package moodle

import (
	"moodle/internal/core/domain"
	"moodle/internal/core/port"
)


type Service struct {
	moodleRepository port.MoodleRepository
	mariaRepository port.MariaRepository
	k8sRepository port.K8sRepository
}

func NewService(moodleRepository port.MoodleRepository, mariaRepository port.MariaRepository, k8sRepository port.K8sRepository) *Service {
	return &Service{
		moodleRepository: moodleRepository,
		mariaRepository: mariaRepository,
		k8sRepository: k8sRepository,
	}
}

func (s *Service) Create(moodle domain.Moodle) (domain.Moodle, error) {
	// create namespace
	// err := s.k8sRepository.CreateNamespace(moodle.Id)
	// if err != nil {
	// 	return domain.Moodle{}, err
	// }
	// create statefulset
	// create service
	// create pvc
	// create db
	// create moodle

	return domain.Moodle{}, nil
}


// get moodle by moodleID from moodles collection in mongodb
func (s *Service) Get(moodleID string) (domain.Moodle, error) {
	moodle, err := s.moodleRepository.Get(moodleID)
	if err != nil {
		return domain.Moodle{}, err
	}
	return moodle, nil
}

// validate if Moodle.WebSiteName is unique
func (s *Service) ValidateSitename(sitename string) bool {
	moodles, err := s.moodleRepository.GetSiteName(sitename)
	if err != nil {
		return false
	}
	if len(moodles) > 0 {
		return false
	}
	return true
}

// List moodles
func (s *Service) List(email string) ([]domain.Moodle, error) {
	moodles, err := s.moodleRepository.GetAll(email)
	if err != nil {
		return nil, err
	}
	return moodles, nil
}


// List courses
func (s *Service) ListCourses() ([]domain.Course, error) {
	courses, err := s.moodleRepository.ListCourses()
	if err != nil {
		return nil, err
	}
	return courses, nil
}


// Delete moodle
func (s *Service) Delete(moodleID string) error {
	err := s.moodleRepository.Delete(moodleID)
	if err != nil {
		return err
	}
	return nil
}


// search moodle by email
func (s *Service) Search(email string) ([]domain.Moodle, error) {
	moodles, err := s.moodleRepository.GetAll(email)
	if err != nil {
		return nil, err
	}
	return moodles, nil
}
