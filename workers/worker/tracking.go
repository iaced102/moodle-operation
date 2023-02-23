package worker

import (
	repo "moodle/internal/repository/moodle"
	mongo "moodle/pkg/mongodbiface"
	"net/http"
)

type TrackingWorker struct {
    mongo *repo.MongoDB
}

func NewTrackingWorker(m mongo.DB) *TrackingWorker{
	return &TrackingWorker{
		mongo: repo.NewMongoDB(m),
	}
}

func (w *TrackingWorker) UpdateMoodleStatus() error {
	// get all moodle from moodletracking collection where isupdate = false
	moodles, err := w.mongo.GetAllMoodleTracking(false)
	if err != nil {
		return err
	}
	for _, moodle := range moodles {
		go w.checkMoodleStatus(moodle.MoodleId, "http://" + moodle.SiteName+".lms.bizflycloud.vn/login/index.php")
	}
	return nil
}

// check status code from website
func (w * TrackingWorker) checkMoodleStatus(moodleid, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	// log.Println("health: ", url, " status code: ", resp.StatusCode)
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		// update status to Online and isupdate to true
		err := w.mongo.UpdateMoodleStatus(moodleid, "Online")
		if err != nil {
			return err
		}
		err = w.mongo.UpdateMoodleTracking(moodleid, true)
		if err != nil {
			return err
		}
	} 
	// log status code for site
	return nil
}

