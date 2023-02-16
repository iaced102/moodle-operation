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
			// clone course table
			log.Println("cloning mdl_course table")
			w.mariaRepo.CloneCourseTable(tracking.DbName, "mdl_course_deleted")
			log.Println("cloned mdl_course table")
			// clone course_categories table
			log.Println("cloning mdl_course_categories table")
			w.mariaRepo.CloneCourseCategoriesTable(tracking.DbName, "mdl_course_categories_deleted")
			log.Println("cloned mdl_course_categories table")
			// delte row from mdl_course table where id from 32-71
			log.Println("deleting row from mdl_course table")
			for i := 32; i <= 71; i++ {
				w.mariaRepo.DeleteRow(tracking.DbName, "mdl_course", i)
			}
			log.Println("deleted row from mdl_course table")
			log.Println("deleting row from mdl_course_categories table")
			for i := 27; i <= 39; i++ {
				w.mariaRepo.DeleteRow(tracking.DbName, "mdl_course_categories", i)
			}
			log.Println("deleted row from mdl_course_categories table")
			}
		// check if sitenameupdate is false then update shortname, fullname to mdl_course
		if tracking.SiteNameUpdate == false {
			log.Println("updating shortname, fullname to mdl_course")
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
			log.Println("updated shortname, fullname to mdl_course")
		}
	}
	return nil
}
