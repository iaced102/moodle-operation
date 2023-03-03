package worker

import (
	"context"
	"database/sql"
	"log"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
	"strings"

	"moodle/internal/core/domain"
)

type SitenameWorker struct {
	mariaRepo *repo.MariaDB
	mongoRepo *repo.MongoDB
}

// new mariadb worker
func NewSitenameWorker(mongo mongodbiface.DB, maria *sql.DB) *SitenameWorker {
	mariaRepo := repo.NewMariaDB(maria)
	mongoRepo := repo.NewMongoDB(mongo)
	
	return &SitenameWorker{
		mariaRepo: mariaRepo,
		mongoRepo: mongoRepo,
	}
}

// restore database then update to maria_tracking collection
func (w *SitenameWorker) UpdateSitename() error {
	// sitename tracking
	siteNameTracking, err := w.mongoRepo.GetAllSiteNameTracking()
	if err != nil {
		log.Println(err)
		return err
	}
	for siteNameTracking.Next(context.Background()) {
		var tracking domain.SitenameTracking
		err := siteNameTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		go w.Update(tracking.DbName, tracking.SiteName)
	}
	return nil
}

func (w *SitenameWorker) Update(dbname, sitename string) error {
	// check if database is mariaTracking
	mariaTracking, err := w.mongoRepo.GetMariaTracking(strings.ReplaceAll(dbname, "_", "-"))
	if err != nil {
		log.Println(err)
		return err
	}
	if mariaTracking.DbStatus == "Restoring" {
		log.Printf("database %s is restoring", dbname)
		return nil
	}
	if mariaTracking.DbStatus == "Restore Failed" {
		log.Printf("database %s is Restore Failed", dbname)
		return nil
	}
	if mariaTracking.DbStatus == "Deleting" {
		return nil
	}
	if mariaTracking.DbStatus != "Created" {
		log.Printf("database %s is not created yet", dbname)
		return nil
	}
	if mariaTracking.DbStatus == "Created" {
		// get pre_installed_course_tracking
		preInstalledCourseTracking, err := w.mongoRepo.GetPreInstalledCourseTracking(strings.ReplaceAll(dbname, "_", "-"))
		if err != nil {
			log.Println(err)
			return err
		}
		if preInstalledCourseTracking.Status == "Updated" {
			log.Println("updating shortname, fullname to mdl_course for " + sitename)
			err = w.mariaRepo.UpdateDB(dbname, sitename, sitename)
			if err != nil {
				log.Println(err)
				return err
			}
			err := w.mongoRepo.UpdateSitenameTracking(domain.SitenameTracking{ 
				MoodleId: strings.ReplaceAll(dbname, "_", "-"),
				DbName: dbname,
				SiteName: sitename,
				IsUpdate: true,
			})
			if err != nil {
				log.Println(err)
				return err
			}
		// update sitename
			err = w.mongoRepo.UpdateSitename(domain.SitenameTracking{
				MoodleId: strings.ReplaceAll(dbname, "_", "-"),
				SiteName: sitename,
			})
			if err != nil {
				return err
			}
			log.Println("updated shortname, fullname to mdl_course for " + sitename)
		}
	}
	return nil
}
