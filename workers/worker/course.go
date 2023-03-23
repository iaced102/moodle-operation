package worker

import (
	"context"
	"database/sql"
	"fmt"
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
			// update status to updating
			tracking.Status = "Updating"
			err = w.mongoRepo.UpdatePreInstalledCourseTracking(tracking)
			// delete cate that in PreCates but not in cates
			// time.Sleep(5 * time.Second)
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
			updatedata := `O:8:"stdClass":17:{s:5:"title";s:28:"CHƯƠNG TRÌNH ĐÀO TẠO ";s:8:"subtitle";s:0:"";s:11:"button_text";s:16:"View All Courses";s:4:"body";s:1:"0";s:5:"items";s:1:"8";s:8:"color_bg";s:18:"rgb(252, 246, 246)";s:11:"color_title";s:17:"rgb(133, 25, 204)";s:14:"color_subtitle";s:7:"#6f7074";s:13:"color_overlay";s:17:"rgba(10,10,10,.5)";s:11:"color_hover";s:18:"rgb(116, 125, 177)";s:9:"color_btn";s:17:"rgb(102, 66, 156)";s:12:"button_bdrrd";s:2:"50";s:10:"categories";a:%d:{%s}s:5:"style";s:1:"1";s:14:"ccn_margin_top";s:1:"0";s:17:"ccn_margin_bottom";s:1:"0";s:13:"ccn_css_class";s:0:"";}`
			result := ""
			for i, cate := range tracking.CourseId {
				s := 0
				if cateCourses.CateCourse[cate-1].CategoryId == 5 {
					s = 1
					result += fmt.Sprintf(`i:%d;s:%d:"%d";`, i, s, cateCourses.CateCourse[cate-1].CategoryId)
				} else {
					s = 2
					result += fmt.Sprintf(`i:%d;s:%d:"%d";`, i, s, cateCourses.CateCourse[cate-1].CategoryId)
				}
			}
			updatedata = fmt.Sprintf(updatedata, len(tracking.CourseId), result)
			updatedata_encode := helpers.Encode(updatedata)
			err = w.mariaRepo.UpdateBlockInstance(strings.ReplaceAll(tracking.MoodleId, "-", "_"), updatedata_encode)
			if err != nil {
				log.Println(err)
				return err
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
			// update status to updated
			tracking.Status = "Updating"
			err = w.mongoRepo.UpdatePreInstalledCourseTracking(tracking)
			for _, v := range cates {
				err = w.mariaRepo.CopyCategoryRow(strings.ReplaceAll(tracking.MoodleId, "-", "_"), cateCourses.CateCourse[v-1])
				if err != nil {
					log.Println(err)
					return err
				}
			// update status to updated
			tracking.Status = "Updated"
			err = w.mongoRepo.UpdatePreInstalledCourseTracking(tracking)
			}
			// get moodle by id
			moodle, err := w.mongoRepo.Get(mariaTracking.MoodleId)
			updatedata := `O:8:"stdClass":17:{s:5:"title";s:28:"CHƯƠNG TRÌNH ĐÀO TẠO ";s:8:"subtitle";s:0:"";s:11:"button_text";s:16:"View All Courses";s:4:"body";s:1:"0";s:5:"items";s:1:"8";s:8:"color_bg";s:18:"rgb(252, 246, 246)";s:11:"color_title";s:17:"rgb(133, 25, 204)";s:14:"color_subtitle";s:7:"#6f7074";s:13:"color_overlay";s:17:"rgba(10,10,10,.5)";s:11:"color_hover";s:18:"rgb(116, 125, 177)";s:9:"color_btn";s:17:"rgb(102, 66, 156)";s:12:"button_bdrrd";s:2:"50";s:10:"categories";a:%d:{%s}s:5:"style";s:1:"1";s:14:"ccn_margin_top";s:1:"0";s:17:"ccn_margin_bottom";s:1:"0";s:13:"ccn_css_class";s:0:"";}`
			result := ""
			for i, cate := range moodle.PreInstalledCourse {
				s := 0
				if cateCourses.CateCourse[cate-1].CategoryId == 5 {
					s = 1
					result += fmt.Sprintf(`i:%d;s:%d:"%d";`, i, s, cateCourses.CateCourse[cate-1].CategoryId)
				} else {
					s = 2
					result += fmt.Sprintf(`i:%d;s:%d:"%d";`, i, s, cateCourses.CateCourse[cate-1].CategoryId)
				}
			}
			updatedata = fmt.Sprintf(updatedata, len(tracking.CourseId), result)
			updatedata_encode := helpers.Encode(updatedata)
			err = w.mariaRepo.UpdateBlockInstance(strings.ReplaceAll(tracking.MoodleId, "-", "_"), updatedata_encode)
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}
