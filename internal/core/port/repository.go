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
	CreateMailTracking(mailTracking domain.Moodle) error
	DeleteMailTracking(mailTracking domain.MailTracking) error
	CreateMariaTracking(mariaTracking domain.Moodle) error
	CreateMoodleConfig(moodleconfig domain.Moodle) error
	UpdateLogo(logotracking domain.LogoTracking) error
	UpdateFavicon(favicontracking domain.FaviconTracking) error
	UpdatePreInstalledCourse(moodle domain.Moodle) error
	UpdateMariaTracking(mariaTracking domain.MariaTracking) error
	GetMariaTracking(moodleid string) (domain.MariaTracking, error)
	CreateMoodleTracking(moodleTracking domain.Moodle) error
	CreateSiteNameTracking(moodle domain.Moodle) error
	UpdateSitename(sitenameTracking domain.SitenameTracking) error
	UpdateSitenameTracking(sitenameTracking domain.SitenameTracking) error
	GetLogoTracking(moodleid string) (domain.LogoTracking, error)
	GetFaviconTracking(moodleid string) (domain.FaviconTracking, error)
	CreateLogoTracking(moodle domain.Moodle, url string) error
	CreateFaviconTracking(moodle domain.Moodle, url string) error

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
	ApplyIngress(namespace, sitename string) error
	DeleteNamespace(namespace string) error
	UpdateIngress(namespace, sitename string) error
}
