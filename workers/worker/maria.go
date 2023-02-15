package worker

import (
	"context"
	"database/sql"
	"log"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"

	"moodle/internal/core/domain"
)

type MariaWorker struct {
	mariaRepo *repo.MariaDB
	mongoRepo *repo.MongoDB
}

// new mariadb worker
func NewMariaWorker(mongo mongodbiface.DB, maria *sql.DB) *MariaWorker {
	mariaRepo := repo.NewMariaDB(maria)
	mongoRepo := repo.NewMongoDB(mongo)
	
	return &MariaWorker{
		mariaRepo: mariaRepo,
		mongoRepo: mongoRepo,
	}
}

// restore database then update to maria_tracking collection
func (w *MariaWorker) Restore() error {
	// get all from maria_tracking
	mariaTracking, err := w.mongoRepo.GetMariaTracking()
	if err != nil {
		log.Println(err)
		return err
	}
	
	// for each maria_tracking then restore database and update to maria_tracking
	for mariaTracking.Next(context.Background()) {
		var tracking domain.MariaTracking
		err := mariaTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		// restore database
		// check if status is Creating
		if tracking.DbStatus == "Creating" {
			err = w.mariaRepo.RestoreDB(tracking.DbName, tracking.FilePath)
			if err != nil {
				log.Println(err)
				return err
			}
			// update DbStatus to maria_tracking
			err = w.mongoRepo.UpdateMariaTracking(tracking)
			if err != nil {
				log.Println(err)
				return err
			}
		}
		// check if sitenameupdate is false then update shortname, fullname to mdl_course
		if tracking.SiteNameUpdate == false {
			err = w.mariaRepo.UpdateDB(tracking.DbName, tracking.SiteName, tracking.SiteName)
			if err != nil {
				log.Println(err)
				return err
			}
			// update sitenameupdate to true
			err = w.mongoRepo.UpdateMariaTrackingSiteNameUpdate(tracking)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}
