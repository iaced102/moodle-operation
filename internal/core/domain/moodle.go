package domain

import (
	"time"

	"github.com/google/uuid"
)

type Moodle struct {
	Id  string `json:"id"`
	Email string `json:"email"`
	Ip string `json:"ip"`
	Name string `json:"name"`
	LbName string `json:"lb_name"`
	WebSiteName string `json:"website_name"`
	PreInstalledCourse []int `json:"pre_installed_course"`
	PackageName string `json:"packages_name"`
	Packages Package `json:"packages"`
	AutoScale bool `json:"autoscale"`
	DocumentsStorageExtra int `json:"documents_storage_extra"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Package struct {
	Type string `json:"type"`
	Price int `json:"price"`
	Name string `json:"name"`
	Ccu int `json:"ccu"`
	AccountMax int `json:"account_max"`
	DocumentStorage int `json:"document_storage"`
	MoodleVersion string `json:"moodle_version"`
	BackupNum int `json:"backup_num"`
	CcuExtraMax int `json:"ccu_extra_max"`
	DocumentStorageExtraMax int `json:"document_storage_extra_max"`
}

type Course struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Content string `json:"content"`
	Thumb_Url string `json:"thumb_url"`
	Highlight []string `json:"highlight"`
	Routine []string `json:"routine"`
	Course_Preview []string `json:"course_preview"`
}

type MoodleConfiguration struct {
	MoodleId string `json:"moodle_id"`
	FaviconUrl string `json:"favicon_url"`
	FaviconPath string `json:"favicon_path"`
	LogoUrl string `json:"logo_url"`
	LogoPath string `json:"logo_path"`
}

type Banner struct {
	ImageUrl string `json:"image_url"`
	Slogan	 string `json:"slogan"`
}

type Vision struct {
	ImageUrl string `json:"image_url"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type PreInstalledCourseTracking struct {
	MoodleId string `json:"moodle_id"`
	CourseId []int `json:"course_id"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CateCourse struct {
	PreInstalledCourseId int `json:"pre_installed_course_id"`
	CategoryId int `json:"category_id"`
	SubCate []int `json:"sub_cate"`
	Courses []int `json:"courses"`
}

type CateCourses struct {
	CateCourse []CateCourse `json:"cate_course"`
}

func NewMoodleID() string {
	return uuid.New().String()
}

