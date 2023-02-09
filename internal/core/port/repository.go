package port

import "moodle/internal/core/domain"

type MoodleRepository interface {
	InsertOne(domain.Moodle) error
	Get(moodleID string) (domain.Moodle, error)
	GetAll(email string) ([]domain.Moodle, error)
	Delete(moodleID string) error
	GetSiteName(siteName string) ([]domain.Moodle, error)
	ListCourses() ([]domain.Course, error)
}

type MariaRepository interface {
	CreateDB(dbname string) error
	RestoreDB(dbname, filepath string) error
	DropDB(dbname string) error
	UpdateDB(dbname, shortname, fullname string) error

}

type K8sRepository interface {
	CreateNamespace(namespace string) error
	ApplyStatefulSet(namespace string) error
	ApplyService(namespace string) error
	ApplyPVC(yaml string) error
	DeleteNamespace(namespace string) error
}
