package domain

import "time"

type Moodle struct {
	Id  string `json:"id"`
	Email string `json:"email"`
	Ip string `json:"ip"`
	Name string `json:"name"`
	LbName string `json:"lb_name"`
	WebSiteName string `json:"website_name"`
	PreInstalledCourse []int `json:"pre_installed_course"`
	Packages MoodlePackages `json:"packages"`
	AutoScale bool `json:"autoscale"`
	DocumentsStorageExtra int `json:"documents_storage_extra"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoodlePackages struct {
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
