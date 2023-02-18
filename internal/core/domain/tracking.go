package domain

import "time"

type MariaTracking struct {
	MoodleId string `json:"moodle_id"`
	DbName string `json:"db_name"`
	DbStatus string `json:"db_status"`
	FilePath string `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SitenameTracking struct {
	MoodleId string `json:"moodle_id"`
	DbName string `json:"db_name"`
	SiteName string `json:"site_name"`
	IsUpdate bool `json:"site_name_update"`
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

type LogoTracking struct {
	MoodleId string `json:"moodle_id"`
	FilePath string `json:"file_path"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FaviconTracking struct {
	MoodleId string `json:"moodle_id"`
	FilePath string `json:"file_path"`
	Status string `json:"favicon_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoodleTracking struct {
	MoodleId string `json:"moodle_id"`
	SiteName string `json:"site_name"`
	IsUpdate bool `json:"is_update"`
	UpdatedAt time.Time `json:"updated_at"`
}
