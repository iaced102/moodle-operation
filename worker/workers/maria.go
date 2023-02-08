package worker

import (
	"context"
	"log"
	mariaAdapter "moodle/pkg/mariadbiface"
	mongoAdapter "moodle/pkg/mongodbiface"
	"moodle/handler"
	"moodle/config"

	"go.mongodb.org/mongo-driver/bson"
)

type MariaWorker struct {
	MariaAdapter *mariaAdapter.MariaAdapter
	MongoAdapter mongoAdapter.MongoDB
}

// new mariadb worker
func NewMariaWorker(adapter mongoAdapter.MongoDB) *MariaWorker {
	mariaclient := mariaAdapter.NewMariaAdapter(config.MARIAHOSTW, 3306, config.MARIAUSER, config.MARIAPASSWORD)
	
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
		var tracking handler.DBTracking
		err := mariaTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		// restore database
		// check if status is Creating
		if tracking.DbStatus == "Creating" {
			err = w.MariaAdapter.RestoreDatabase(tracking.DbName, tracking.FilePath)
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
				err = w.MariaAdapter.Update(tracking.DbName, tracking.SiteName, tracking.SiteName)
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
