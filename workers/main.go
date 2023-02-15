package main

import (
	"context"
	"database/sql"
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
	maria := NewMariaDB()
	go LBWorker(mongo)
	go SendmailWorker(mongo)
	go TrackingWorker(mongo)
	go MariaWorker(mongo, maria)
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

func MariaWorker(mongo mongodbiface.DB, maria *sql.DB) {
	log.Println("Starting Maria worker")
	for {
		time.Sleep(5 * time.Second)
		worker.NewMariaWorker(mongo, maria).Restore()
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

func NewMariaDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.MARIAUSER, config.MARIAPASSWORD, config.MARIAHOSTW, config.MARIAPORT, "moodle"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	fmt.Println("Connected to MariaDB!")
	return db

}


// run worker forever
func main() {
	mongo := NewMongoDB()
	Start(mongo)
	select {}
}
