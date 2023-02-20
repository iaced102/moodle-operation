package port

import (
	"mime/multipart"
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
	UpdateLogo(moodleID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError)
	UpdateFavicon(moodleID string, file multipart.File, header *multipart.FileHeader) (map[string]string, *apperrors.AppError)
	UpdatePreInstalledCourse(moodle domain.Moodle) (map[string]string, *apperrors.AppError)
	UpdateSiteName(moodle domain.SitenameTracking) (map[string]string, *apperrors.AppError)
	GetLogo(moodleID string) (map[string]string, *apperrors.AppError)
	GetFavicon(moodleID string) (map[string]string, *apperrors.AppError)
}

