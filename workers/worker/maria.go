package worker

import (
	"context"
	"database/sql"
	"log"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"

	"moodle/internal/core/domain"
)

type MariaWroker struct {
	mariaRepo *repo.MariaDB
	mongoRepo *repo.MongoDB
}

// new mariadb worker
func NewMariaWorker(mongo mongodbiface.DB, maria *sql.DB) *MariaWroker {
	mariaRepo := repo.NewMariaDB(maria)
	mongoRepo := repo.NewMongoDB(mongo)
	
	return &MariaWroker{
		mariaRepo: mariaRepo,
		mongoRepo: mongoRepo,
	}
}

// restore database then update to maria_tracking collection
func (w *MariaWroker) Restore() error {
	// get all from maria_tracking
	mariaTracking, err := w.mongoRepo.GetAllMariaTracking()
	if err != nil {
		log.Println(err)
		return err
	}
	for mariaTracking.Next(context.Background()) {
		var tracking domain.MariaTracking
		err := mariaTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		go w.RestoreDB(tracking.DbName, tracking.FilePath)
	}

	// // clone course table
	// log.Println("cloning mdl_course table")
	// w.mariaRepo.CloneCourseTable(tracking.DbName, "mdl_course_deleted")
	// log.Println("cloned mdl_course table")
	// // clone course_categories table
	// log.Println("cloning mdl_course_categories table")
	// w.mariaRepo.CloneCourseCategoriesTable(tracking.DbName, "mdl_course_categories_deleted")
	// log.Println("cloned mdl_course_categories table")
	// // delte row from mdl_course table where id from 32-71
	// log.Println("deleting row from mdl_course table")
	// for i := 32; i <= 71; i++ {
	// 	w.mariaRepo.DeleteRow(tracking.DbName, "mdl_course", i)
	// }
	// log.Println("deleted row from mdl_course table")
	// log.Println("deleting row from mdl_course_categories table")
	// for i := 27; i <= 39; i++ {
	// 	w.mariaRepo.DeleteRow(tracking.DbName, "mdl_course_categories", i)
	// }
	// log.Println("deleted row from mdl_course_categories table")
	return nil
}

// restore database 
func (w *MariaWroker) RestoreDB(dbname, filepath string) error {
	// update DbStatus to maria_tracking
	mariatracking := domain.MariaTracking{
		DbName: dbname,
		DbStatus: "Restoring",
	}
	err := w.mongoRepo.UpdateMariaTracking(mariatracking)
	if err != nil {
		log.Println(err)
		return err
	}
	err = w.mariaRepo.RestoreDB(dbname, filepath)
	if err != nil {
		log.Println(err)
		tracking := domain.MariaTracking{
			DbName: dbname,
			DbStatus: "Restore Failed",
		}
		err := w.mongoRepo.UpdateMariaTracking(tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	mariatracking = domain.MariaTracking{
		DbName: dbname,
		DbStatus: "Created",
	}
	err = w.mongoRepo.UpdateMariaTracking(mariatracking)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
