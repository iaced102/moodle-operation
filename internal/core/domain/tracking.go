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

type MariaTracking struct {
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
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VisionImageTracking struct {
	MoodleId string `json:"moodle_id"`
	FilePath string `json:"file_path"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerImageTracking struct {
	MoodleId string `json:"moodle_id"`
	BannerId string `json:"banner_id"`
	FilePath string `json:"file_path"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BannerSloganTracking struct {
	MoodleId string `json:"moodle_id"`
	SloganId string `json:"slogan_id"`
	Slogan string `json:"slogan"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoURLTracking struct {
	MoodleId string `json:"moodle_id"`
	VideoURL string `json:"video_url"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VisionContentTracking struct {
	MoodleId string `json:"moodle_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status string `json:"logo_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
