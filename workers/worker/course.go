package worker

import (
	"context"
	"database/sql"
	"log"
	repo "moodle/internal/repository/moodle"
	helpers "moodle/pkg/helpers"
	"moodle/pkg/mongodbiface"
	"strings"

	"moodle/internal/core/domain"
)

type CourseWorker struct {
	mariaRepo *repo.MariaDB
	mongoRepo *repo.MongoDB
}

// new course worker
func NewCourseWorker(mongo mongodbiface.DB, maria *sql.DB) *CourseWorker {
	mariaRepo := repo.NewMariaDB(maria)
	mongoRepo := repo.NewMongoDB(mongo)
	
	return &CourseWorker{
		mariaRepo: mariaRepo,
		mongoRepo: mongoRepo,
	}
}


func (w *CourseWorker) Update(){
	go w.Delete()
	go w.Add()
}

// delete course that do not use when create
func (w *CourseWorker) Delete() error {
	// get all from pre_installed_course_tracking where status is Creating
	courseTracking, err := w.mongoRepo.ListPreInstalledCourseTracking("Creating")
	if err != nil {
		log.Println(err)
		return err
	}
	PreCates := []int{1,2,3,4,5,6,7,8}
	cateCourses , err := w.mongoRepo.GetCateCourse()
	if err != nil {
		log.Println(err)
		return err
	}
	for courseTracking.Next(context.Background()) {
		var tracking domain.PreInstalledCourseTracking
		err := courseTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		cates := tracking.CourseId
		// check if dbstatus is Created  from maria_tracking
		mariaTracking, err := w.mongoRepo.GetMariaTracking(tracking.MoodleId)
		if mariaTracking.DbStatus == "Created" {
			// delete cate that in PreCates but not in cates
			for _, v := range PreCates {
				if !helpers.Contains(cates, v) {
					err = w.mariaRepo.DeleteCategory(strings.ReplaceAll(tracking.MoodleId, "-", "_"), cateCourses.CateCourse[v-1])
					if err != nil {
						log.Println(err)
						return err
					}
				}
			// update status to updated
			tracking.Status = "Updated"
			err = w.mongoRepo.UpdatePreInstalledCourseTracking(tracking)
			}
		}
	}
	return nil
}


// update course if that was added
func (w *CourseWorker) Add() error {
	// get all from pre_installed_course_tracking where status is Creating
	courseTracking, err := w.mongoRepo.ListPreInstalledCourseTracking("Pending")
	if err != nil {
		log.Println(err)
		return err
	}
	PreCates := []int{1,2,3,4,5,6,7,8}
	cateCourses , err := w.mongoRepo.GetCateCourse()
	if err != nil {
		log.Println(err)
		return err
	}
	for courseTracking.Next(context.Background()) {
		var tracking domain.PreInstalledCourseTracking
		err := courseTracking.Decode(&tracking)
		if err != nil {
			log.Println(err)
			return err
		}
		cates := tracking.CourseId
		// check if dbstatus is Created  from maria_tracking
		mariaTracking, err := w.mongoRepo.GetMariaTracking(tracking.MoodleId)
		if mariaTracking.DbStatus == "Created" {
			// delete cate that in PreCates but not in cates
			for _, v := range PreCates {
				if helpers.Contains(cates, v) {
					err = w.mariaRepo.CopyCategoryRow(strings.ReplaceAll(tracking.MoodleId, "-", "_"), cateCourses.CateCourse[v-1])
					if err != nil {
						log.Println(err)
						return err
					}
				}
			// update status to updated
			tracking.Status = "Updated"
			err = w.mongoRepo.UpdatePreInstalledCourseTracking(tracking)
			}
		}
	}
	return nil
}
