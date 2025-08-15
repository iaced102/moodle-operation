package domain

import "time"

type MoodleMariaMapping struct {
	MoodleId  string    `json:"moodle_id" bson:"moodle_id"`
	MariaId   string    `json:"maria_id" bson:"maria_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type MariaConfig struct {
	ID        string    `json:"id" bson:"id"`
	DBHost    string    `json:"dbhost" bson:"dbhost"`
	DBPort    int       `json:"dbport" bson:"dbport"`
	DBUser    string    `json:"dbuser" bson:"dbuser"`
	DBPass    string    `json:"dbpass" bson:"dbpass"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
