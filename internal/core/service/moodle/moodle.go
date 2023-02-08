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

func NewService(moodleRepository port.MoodleRepository, mariaRepository port.MariaRepository) *Service {
	return &Service{
		moodleRepository: moodleRepository,
		mariaRepository: mariaRepository,
	}
}

func (s *Service) Create(moodle domain.Moodle) (domain.Moodle, error) {
	// create namespace
	err := s.k8sRepository.CreateNamespace(moodle.Id)
	if err != nil {
		return domain.Moodle{}, err
	}
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
