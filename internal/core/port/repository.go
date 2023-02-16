package port

import "moodle/internal/core/domain"

type MongoRepository interface {
	Create(moodle domain.Moodle) error
	GetAll(email string) ([]domain.Moodle, error)
	Get(moodleID string) (domain.Moodle, error)
	Search(email, name string) ([]domain.Moodle, error)
	ListCourses() ([]domain.Course, error)
	ListPackages() ([]domain.Package, error)
	Delete(moodleID string) error
	GetSiteName(name string) ([]domain.Moodle, error)
	GetPackage(name string) (domain.Package, error)
	CreateLbTracking(lbTracking domain.Moodle) error
	CreateMailTracking(mailTracking domain.Moodle) error
	CreateMariaTracking(mariaTracking domain.Moodle) error
	CreateMoodleConfig(moodleconfig domain.Moodle) error
	UpdateLogo(logotracking domain.LogoTracking) error
	UpdateFavicon(favicontracking domain.FaviconTracking) error
	UpdateVisionImage(visionimagetracking domain.VisionImageTracking) error
	UpdateBannerImage(bannerimagetracking domain.BannerImageTracking) error
	GetBannerImage(moodleid string) (domain.BannerImageTracking, error)
	UpdateVideoURL(videourltracking domain.VideoURLTracking) error
	UpdateVisionContent(visioncontenttracking domain.VisionContentTracking) error
	UpdateBannerSlogan(bannerslogantracking domain.BannerSloganTracking) error
	UpdatePreInstalledCourse(moodle domain.Moodle) error
}

type MariaRepository interface {
	CreateDB(dbname string) error
	RestoreDB(dbname, filepath string) error
	DropDB(dbname string) error
	UpdateDB(dbname, shortname, fullname string) error

}

type K8sRepository interface {
	ListNamespaces() ([]string, error)
	CreateNamespace(namespace string) error
	ApplyStatefulSet(namespace string) error
	ApplyService(namespace string) error
	ApplyPVC(namespace string) error
	DeleteNamespace(namespace string) error
}
