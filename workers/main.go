package workers

import (
	"moodle/pkg/mongodbiface"
	"moodle/workers/worker"
	"time"

	"log"
)

func Start(mongo mongodbiface.DB) {
	for {
		time.Sleep(5 * time.Second)
		log.Println("Starting LB worker")
		worker.NewLBWorker(mongo).GetLBInstances()
	}
}
