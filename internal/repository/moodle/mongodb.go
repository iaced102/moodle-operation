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
	_, err := m.db.Collection("moodles").DeleteOne(context.Background(), bson.M{"id": moodleID})
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

// create moodle tracking
func (m *MongoDB) CreateMoodleTracking(moodle domain.Moodle) error {
	var moodleTracking domain.MoodleTracking
	moodleTracking.MoodleId = moodle.Id
	moodleTracking.SiteName = moodle.Name
	moodleTracking.IsUpdate = false
	moodleTracking.UpdatedAt = time.Now()
	_, err := m.db.Collection("moodle_tracking").InsertOne(context.Background(), moodleTracking)
	if err != nil {
		return err
	}
	return nil
}

// update moodle status 
func (m *MongoDB) UpdateMoodleStatus(moodleid, status string) error {
	_, err := m.db.Collection("moodles").UpdateOne(context.Background(), bson.M{"id": moodleid}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	return nil
}

// update moodle tracking
func (m *MongoDB) UpdateMoodleTracking(moodleid string, isupdate bool) error {
	_, err := m.db.Collection("moodle_tracking").UpdateOne(context.Background(), bson.M{"moodleid": moodleid}, bson.M{"$set": bson.M{"isupdate": isupdate, "updatedat": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// get all moodle tracking
func (m *MongoDB) GetAllMoodleTracking(isupdate bool) ([]domain.MoodleTracking, error) {
	var moodleTrackings []domain.MoodleTracking
	cursor, err := m.db.Collection("moodle_tracking").Find(context.Background(), bson.M{"isupdate": isupdate})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var moodleTracking domain.MoodleTracking
		err := cursor.Decode(&moodleTracking)
		if err != nil {
			return nil, err
		}
		moodleTrackings = append(moodleTrackings, moodleTracking)
	}
	return moodleTrackings, nil
}

// create maria tracking
func (m *MongoDB) CreateMariaTracking(moodle domain.Moodle) error {
	var mariaTracking domain.MariaTracking
	mariaTracking.MoodleId = moodle.Id
	mariaTracking.DbName = strings.ReplaceAll(moodle.Id, "-", "_")
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

// Get all maria tracking where dbstatus is creating
func (m *MongoDB) GetAllMariaTracking() (*mongo.Cursor, error) {
	cursor, err := m.db.Collection("maria_tracking").Find(context.Background(), bson.M{"dbstatus": "Creating"})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// create sitename tracking
func (m *MongoDB) CreateSiteNameTracking(moodle domain.Moodle) error {
	var sitenameTracking domain.SitenameTracking
	sitenameTracking.MoodleId = moodle.Id
	sitenameTracking.DbName = strings.ReplaceAll(moodle.Id, "-", "_")
	sitenameTracking.SiteName = moodle.Name
	sitenameTracking.IsUpdate = false
	sitenameTracking.CreatedAt = time.Now()
	sitenameTracking.UpdatedAt = time.Now()
	_, err := m.db.Collection("sitename_tracking").InsertOne(context.Background(), sitenameTracking)
	if err != nil {
		return err
	}
	return nil
}

// get all sitename tracking where isupdate is false
func (m *MongoDB) GetAllSiteNameTracking() (*mongo.Cursor, error) {
	cursor, err := m.db.Collection("sitename_tracking").Find(context.Background(), bson.M{"isupdate": false})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// update maria tracking filter by dbname and update dbstatus
func (m *MongoDB) UpdateMariaTracking(mariaTracking domain.MariaTracking) error {
	_, err := m.db.Collection("maria_tracking").UpdateOne(context.Background(), bson.M{"dbname": mariaTracking.DbName}, bson.M{"$set": bson.M{"dbstatus": mariaTracking.DbStatus}})
	if err != nil {
		return err
	}
	return nil
}

// update site_name_update filter by dbname and update site_name_update
func (m *MongoDB) UpdateSitenameTracking(mariaTracking domain.SitenameTracking) error {
	_, err := m.db.Collection("sitename_tracking").UpdateOne(context.Background(), bson.M{"moodleid": mariaTracking.MoodleId}, bson.M{"$set": mariaTracking})
	if err != nil {
		return err
	}
	return nil
}

// update sitename
func (m *MongoDB) UpdateSitename(sitenametracking domain.SitenameTracking) error {
	_, err := m.db.Collection("moodles").UpdateOne(context.Background(), bson.M{"id": sitenametracking.MoodleId}, bson.M{"$set": bson.M{"name": sitenametracking.SiteName, "websitename": sitenametracking.SiteName + ".lms.bizflycloud.vn", "updatedat": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// update logotracking
func (m *MongoDB) UpdateLogo(logotracking domain.LogoTracking) error {
	// check if logo exist then update else insert
	var logo domain.LogoTracking
	err := m.db.Collection("logo_tracking").FindOne(context.Background(), bson.M{"moodleid": logotracking.MoodleId}).Decode(&logo)
	if err != nil {
		_, err := m.db.Collection("logo_tracking").InsertOne(context.Background(), logotracking)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = m.db.Collection("logo_tracking").UpdateOne(context.Background(), bson.M{"moodleid": logotracking.MoodleId}, bson.M{"$set": logotracking})
	if err != nil {
		return err
	}
	return nil
}

// update favicontracking
func (m *MongoDB) UpdateFavicon(favicontracking domain.FaviconTracking) error {
	// check if favicon exist then update else insert
	var favicon domain.FaviconTracking
	err := m.db.Collection("favicon_tracking").FindOne(context.Background(), bson.M{"moodleid": favicontracking.MoodleId}).Decode(&favicon)
	if err != nil {
		_, err := m.db.Collection("favicon_tracking").InsertOne(context.Background(), favicontracking)
		if err != nil {
			return err
		}
		return nil
	}
	_, err = m.db.Collection("favicon_tracking").UpdateOne(context.Background(), bson.M{"moodleid": favicontracking.MoodleId}, bson.M{"$set": favicontracking})
	if err != nil {
		return err
	}
	return nil
}

// update pre_installed_course
func (m *MongoDB) UpdatePreInstalledCourse(moodle domain.Moodle) error {
	currentMoodle, err := m.Get(moodle.Id)
	if err != nil {
		return err
	}
	currentMoodle.PreInstalledCourse = moodle.PreInstalledCourse
	// update pre_installed_course in moodle
	_, err = m.db.Collection("moodles").UpdateOne(context.Background(), bson.M{"id": moodle.Id}, bson.M{"$set": currentMoodle})
	if err != nil {
		return err
	}
	return nil
}

// get maria_tracking
func (m *MongoDB) GetMariaTracking(moodleid string) (domain.MariaTracking, error) {
	var mariatracking domain.MariaTracking
	err := m.db.Collection("maria_tracking").FindOne(context.Background(), bson.M{"moodleid": moodleid}).Decode(&mariatracking)
	if err != nil {
		return mariatracking, err
	}
	return mariatracking, nil
}

// get filepath
func (m *MongoDB) GetFilePath(moodleid, repo string) (string, error) {
	var tracking domain.LogoTracking
	err := m.db.Collection(repo).FindOne(context.Background(), bson.M{"moodleid": moodleid}).Decode(&tracking)
	if err != nil {
		return "", err
	}
	return tracking.FilePath, nil
}

// get all logo_tracking
func (m *MongoDB) GetAllLogoTracking() (*mongo.Cursor, error) {
	cursor, err := m.db.Collection("logo_tracking").Find(context.Background(), bson.M{"status": "pending"})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// get all favicon_tracking
func (m *MongoDB) GetAllFaviconTracking() (*mongo.Cursor, error) {
	cursor, err := m.db.Collection("favicon_tracking").Find(context.Background(), bson.M{"status": "Pending"})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// update logo_tracking
func (m *MongoDB) UpdateLogoTracking(logoTracking domain.LogoTracking) error {
	_, err := m.db.Collection("logo_tracking").UpdateOne(context.Background(), bson.M{"moodleid": logoTracking.MoodleId}, bson.M{"$set": bson.M{"status": logoTracking.Status}})
	if err != nil {
		return err
	}
	return nil
}
