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
	Type string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LogoTracking struct {
	MoodleId string `json:"moodle_id"`
	FilePath string `json:"file_path"`
	URL string `json:"url"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FaviconTracking struct {
	MoodleId string `json:"moodle_id"`
	FilePath string `json:"file_path"`
	Status string `json:"favicon_status"`
	URL string `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MoodleTracking struct {
	MoodleId string `json:"moodle_id"`
	SiteName string `json:"site_name"`
	IsUpdate bool `json:"is_update"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NFSTracking struct {
	Path string `json:"path"`
	Size int `json:"size"`
}

type NFSTrackings struct {
	Id int `json:"id"`
	Data []NFSTracking `json:"data"`
}


type StorageAlarmTracking struct {
	MoodleId  string    `json:"moodle_id" bson:"moodle_id"`
	Week      int       `json:"week" bson:"week"`
	Year      int       `json:"year" bson:"year"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type ExtendStorage struct {
	MoodleId      string `json:"moodleId" bson:"moodleId"`
	ExtendStorage string `json:"extend_storage" bson:"extend_storage"`
}

// CCUAlertTracking tracks CCU alerts to prevent spam
type CCUAlertTracking struct {
	MoodleId   string    `json:"moodle_id" bson:"moodle_id"`
	AlertType  string    `json:"alert_type" bson:"alert_type"` // "warn" or "alert"
	CurrentCCU int       `json:"current_ccu" bson:"current_ccu"`
	DefaultCCU int       `json:"default_ccu" bson:"default_ccu"`
	MaxCCU     int       `json:"max_ccu" bson:"max_ccu"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}
