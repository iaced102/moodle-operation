package main

import (
	"context"
	"fmt"
	"log"
	"moodle/config"
	"moodle/pkg/mongodbiface"
	"moodle/workers/worker"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewMongoDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	return client.Database(config.DBNAME)
}

// run worker forever
func main() {
	mongo := NewMongoDB()
	Start(mongo)
	select {}
}
