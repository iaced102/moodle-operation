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

// delete course
func (w *CourseWorker) Update() error {
	// get all from pre_installed_course_tracking where status is Pending
	courseTracking, err := w.mongoRepo.ListPreInstalledCourseTracking()
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
		go w.UpdateCate(strings.ReplaceAll(tracking.MoodleId, "-", "_"), tracking.CourseId)
		go w.UpdateCourse(strings.ReplaceAll(tracking.MoodleId, "-", "_"), tracking.CourseId)
	}

	return nil
}

// update category
func (w *CourseWorker) UpdateCate(dbname string, courseid []int) error {
	// cope from mdl_course_categories_backup to mdl_course_categories
	err := w.mariaRepo.CopyCategoryRow(dbname, courseid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// update course
func (w *CourseWorker) UpdateCourse(dbname string, courseid []int) error {
	// cope from mdl_course_backup to mdl_course
	err := w.mariaRepo.CopyCourseRow(dbname, courseid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// delete course that do not use when create
func (w *CourseWorker) Delete() error {
	// get all from pre_installed_course_tracking where status is Creating
	courseTracking, err := w.mongoRepo.ListPreInstalledCourseTracking()
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
		// delete cate that in PreCates but not in cates
		for _, v := range PreCates {
			if !helpers.Contains(cates, v) {
				go w.mariaRepo.DeleteCategory(strings.ReplaceAll(tracking.MoodleId, "-", "_"), cateCourses[v].CategoryId)
			}
		}
	}


	return nil
}
