package moodle

import (
	"context"
	"moodle/internal/core/domain"
	"moodle/pkg/mongodbiface"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	db mongodbiface.DB
}

func NewMongoDB(mongodb mongodbiface.DB) *MongoDB {
    return &MongoDB{db: mongodb,}
}

// get moodle by id
func (m *MongoDB) Get(moodleID string) (domain.Moodle, error) {
	var moodle domain.Moodle
	err := m.db.Collection("moodles").FindOne(context.Background(), bson.M{"id": moodleID}).Decode(&moodle)
	if err != nil {
		return moodle, err
	}
	return moodle, nil
}

// get all moodles
func (m *MongoDB) GetAll(email string) ([]domain.Moodle, error) {
	var moodles []domain.Moodle
	cursor, err := m.db.Collection("moodles").Find(context.Background(), bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var moodle domain.Moodle
		err := cursor.Decode(&moodle)
		if err != nil {
			return nil, err
		}
		moodles = append(moodles, moodle)
	}
	return moodles, nil
}

// delete moodle
func (m *MongoDB) Delete(moodleID string) error {
	_, err := m.db.Collection("moodles").DeleteOne(context.Background(), bson.M{"moodleid": moodleID})
	if err != nil {
		return err
	}
	return nil
}

// get all sitename
func (m *MongoDB) GetSiteName(name string) ([]domain.Moodle, error) {
	var moodles []domain.Moodle
	cursor, err := m.db.Collection("moodles").Find(context.Background(), bson.M{"name": name})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var moodle domain.Moodle
		err := cursor.Decode(&moodle)
		if err != nil {
			return nil, err
		}
		moodles = append(moodles, moodle)
	}
	return moodles, nil
}

// list courses
func (m *MongoDB) ListCourses() ([]domain.Course, error) {
	var courses []domain.Course
	cursor, err := m.db.Collection("courses").Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var course domain.Course
		err := cursor.Decode(&course)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

// list packages
func (m *MongoDB) ListPackages() ([]domain.Package, error) {
	var packages []domain.Package
	cursor, err := m.db.Collection("packages").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var pack domain.Package
		err := cursor.Decode(&pack)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pack)
	}
	return packages, nil
}

// search moodle
func (m *MongoDB) Search(email, name string) ([]domain.Moodle, error) {
	var moodles []domain.Moodle
	cursor, err := m.db.Collection("moodles").Find(context.Background(), bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var moodle domain.Moodle
		err := cursor.Decode(&moodle)
		if err != nil {
			return nil, err
		}
		moodles = append(moodles, moodle)
	}
	return moodles, nil
}

// get package by name
func (m *MongoDB) GetPackage(name string) (domain.Package, error) {
	var pack domain.Package
	err := m.db.Collection("packages").FindOne(context.Background(), bson.M{"Name": name}).Decode(&pack)
	if err != nil {
		return pack, err
	}
	return pack, nil
}

// create moodle
func (m *MongoDB) Create(moodle domain.Moodle) error {
	_, err := m.db.Collection("moodles").InsertOne(context.Background(), moodle)
	if err != nil {
		return err
	}
	return nil
}

// create lbtracking
func (m *MongoDB) CreateLbTracking(moodle domain.Moodle) error {
	var lbTracking domain.LBTracking
	lbTracking.MoodleId = moodle.Id
	lbTracking.LbName = moodle.LbName
	lbTracking.VipAddress = ""
	lbTracking.LbStatus = "Creating"
	lbTracking.CreatedAt = time.Now()
	lbTracking.UpdatedAt = time.Now()
	_, err := m.db.Collection("lb_tracking").InsertOne(context.Background(), lbTracking)
	if err != nil {
		return err
	}
	return nil
}

// create maria tracking
func (m *MongoDB) CreateMariaTracking(moodle domain.Moodle) error {
	var mariaTracking domain.MariaTracking
	mariaTracking.MoodleId = moodle.Id
	mariaTracking.DbName = strings.ReplaceAll(moodle.Id, "-", "_")
	mariaTracking.SiteName = moodle.Name
	mariaTracking.SiteNameUpdate = false
	mariaTracking.DbStatus = "Creating"
	mariaTracking.FilePath = "$HOME/gits/moodle-operator/docker/moodle_seded.sql"
	mariaTracking.CreatedAt = time.Now()
	mariaTracking.UpdatedAt = time.Now()
	_, err := m.db.Collection("maria_tracking").InsertOne(context.Background(), mariaTracking)
	if err != nil {
		return err
	}
	return nil
}

// create mail tracking
func (m *MongoDB) CreateMailTracking(moodle domain.Moodle) error {
	var mailTracking domain.MailTracking
	mailTracking.MoodleId = moodle.Id
	mailTracking.Email = moodle.Email
	mailTracking.IsSent = false
	mailTracking.CreatedAt = time.Now()
	mailTracking.UpdatedAt = time.Now()
	_, err := m.db.Collection("mail_tracking").InsertOne(context.Background(), mailTracking)
	if err != nil {
		return err
	}
	return nil
}

// create moodle config
func (m *MongoDB) CreateMoodleConfig(moodle domain.Moodle) error {
	var moodleConfig domain.MoodleConfiguration
	moodleConfig.MoodleId = moodle.Id
	_, err := m.db.Collection("moodle_configs").InsertOne(context.Background(), moodleConfig)
	if err != nil {
		return err
	}
	return nil
}

// Get all maria tracking return cursor
func (m *MongoDB) GetMariaTracking() (*mongo.Cursor, error) {
	cursor, err := m.db.Collection("maria_tracking").Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// update maria tracking filter by dbname and update dbstatus
func (m *MongoDB) UpdateMariaTracking(mariaTracking domain.MariaTracking) error {
	_, err := m.db.Collection("maria_tracking").UpdateOne(context.Background(), bson.M{"dbname": mariaTracking.DbName}, bson.M{"$set": bson.M{"dbstatus": "Online"}})
	if err != nil {
		return err
	}
	return nil
}

// update site_name_update filter by dbname and update site_name_update
func (m *MongoDB) UpdateMariaTrackingSiteNameUpdate(mariaTracking domain.MariaTracking) error {
	_, err := m.db.Collection("maria_tracking").UpdateOne(context.Background(), bson.M{"dbname": mariaTracking.DbName}, bson.M{"$set": bson.M{"sitenameupdate": true}})
	if err != nil {
		return err
	}
	return nil
}
