package port

import "moodle/internal/core/domain"

type MoodleService interface {
	Create(moodle domain.Moodle) (domain.Moodle, error)
	Get(moodleID string) (domain.Moodle, error)
	List(email string) ([]domain.Moodle, error)
	Search(email string) ([]domain.Moodle, error)
	Delete(moodleID string) error
}

