package workers

import (
	"log"
	"moodle/pkg/mongodbiface"
	"moodle/workers/worker"
	"time"
)

func Start(mongo mongodbiface.DB) {
	go LBWorker(mongo)
	go SendmailWorker(mongo)
	go TrackingWorker(mongo)
}

func LBWorker(mongo mongodbiface.DB) {
	log.Println("Starting LB worker")
	for {
		time.Sleep(5 * time.Second)
		worker.NewLBWorker(mongo).GetLBInstances()
	}
}

func SendmailWorker(mongo mongodbiface.DB) {
	log.Println("Starting Sendmail worker")
	for {
		time.Sleep(5 * time.Second)
		worker.NewSendmailWorker(mongo).SendMailWorkerPool()
	}
}

func TrackingWorker(mongo mongodbiface.DB) {
	log.Println("Starting Tracking worker")
	for {
		time.Sleep(5 * time.Second)
		worker.NewTrackingWorker(mongo).UpdateLbStatus()
	}
}
