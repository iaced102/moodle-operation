package domain

import "time"

type LBTracking struct {
	MoodleId string `json:"moodle_id"`
	LbName string `json:"lb_name"`
	VipAddress string `json:"vip_address"`
	LbStatus string `json:"lb_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DBTracking struct {
	MoodleId string `json:"moodle_id"`
	DbName string `json:"db_name"`
	SiteName string `json:"site_name"`
	SiteNameUpdate bool `json:"site_name_update"`
	DbStatus string `json:"lb_status"`
	FilePath string `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MailTracking struct {
	MoodleId string `json:"moodle_id"`
	Email string `json:"email"`
	IsSent bool `json:"is_sent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

