package port

import (
	"moodle/internal/core/domain"
	"moodle/pkg/apperrors"
)

type MoodleService interface {
	Create(moodle domain.Moodle) (domain.Moodle, *apperrors.AppError)
	List(email string) ([]domain.Moodle, error)
	Get(moodleID string) (domain.Moodle, error)
	Search(email, name string) ([]domain.Moodle, error)
	ListCourses() ([]domain.Course, error)
	ListPackages() ([]domain.Package, error)
	Delete(moodleID string) error
}

