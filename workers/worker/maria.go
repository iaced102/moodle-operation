package worker

import (
	"context"
	"log"
	mariaAdapter "moodle/internal/repository/moodle"
	mongoAdapter "moodle/pkg/mongodbiface"
	"moodle/config"

	"go.mongodb.org/mongo-driver/bson"
)

type DBTracking struct {
	MoodleId       string `bson:"moodleid"`
	DbName         string `bson:"dbname"`
	SiteName       string `bson:"sitename"`
	SiteNameUpdate bool   `bson:"sitenameupdate"`
	DbStatus       string `bson:"dbstatus"`
	FilePath       string `bson:"filepath"`
	CreatedAt      string `bson:"createdat"`
	UpdateAt       string `bson:"updateat"`
}

type MariaWorker struct {
	MariaAdapter *mariaAdapter.MariaDB
	MongoAdapter mongoAdapter.DB
}

// new mariadb worker
func NewMariaWorker(adapter mongoAdapter.DB ) *MariaWorker {
	mariaclient := mariaAdapter.NewMariaDB(config.MARIAHOSTW, 3306, config.MARIAUSER, config.MARIAPASSWORD)
	
	return &MariaWorker{
		MariaAdapter: mariaclient,
		MongoAdapter: adapter,
	}
}

// restore database then update to maria_tracking collection
func (w *MariaWorker) Restore() error {
	mariaTrackingCollection := w.MongoAdapter.Collection("maria_tracking")
	// get all from maria_tracking
	mariaTracking, err := mariaTrackingCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		return err
	}
	
	// for each maria_tracking then restore database and update to maria_tracking
	for mariaTracking.Next(context.Background()) {
		var tracking DBTracking
		err := mariaTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		// restore database
		// check if status is Creating
		if tracking.DbStatus == "Creating" {
			err = w.MariaAdapter.RestoreDB(tracking.DbName, tracking.FilePath)
			if err != nil {
				log.Println(err)
				return err
			}
			// update DbStatus to maria_tracking
			_, err = mariaTrackingCollection.UpdateOne(context.Background(), bson.M{"dbname": tracking.DbName}, bson.M{"$set": bson.M{"dbstatus": "Online"}})
			if err != nil {
				log.Println(err)
				return err
			}
			// check if sitenameupdate is false then update shortname, fullname to mdl_course
			if tracking.SiteNameUpdate == false {
				err = w.MariaAdapter.UpdateDB(tracking.DbName, tracking.SiteName, tracking.SiteName)
				if err != nil {
					log.Println(err)
					return err
				}
				// update sitenameupdate to true
				_, err = mariaTrackingCollection.UpdateOne(context.Background(), bson.M{"dbname": tracking.DbName}, bson.M{"$set": bson.M{"sitenameupdate": true}})
				if err != nil {
					log.Println(err)
					return err
				}
			}
		}
	}
	return nil
}
