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
	go MariaWorker(mongo, maria)
	go CourseWorker(mongo, maria)
	go SitenameWorker(mongo, maria)
	go NFSWorker(mongo)
	go SendmailCreateWorker(mongo, maria)
	go SendmailDeleteWorker(mongo, maria)
	go TrackingWorker(mongo)
}

func SendmailCreateWorker(mongo mongodbiface.DB, maria *sql.DB) {
	log.Println("Starting create Sendmail worker")
	for {
		time.Sleep(3 * time.Second)
		worker.NewSendmailWorker(mongo, maria).SendMailCreateWorkerPool()
	}
}

func SendmailDeleteWorker(mongo mongodbiface.DB, maria *sql.DB) {
	log.Println("Starting delete Sendmail worker")
	for {
		time.Sleep(3 * time.Second)
		worker.NewSendmailWorker(mongo, maria).SendMailDeleteWorkerPool()
	}
}

func TrackingWorker(mongo mongodbiface.DB) {
	log.Println("Starting Tracking worker")
	for {
		time.Sleep(3 * time.Second)
		worker.NewTrackingWorker(mongo).UpdateMoodleStatus()
	}
}

func MariaWorker(mongo mongodbiface.DB, maria *sql.DB) {
	log.Println("Starting Maria worker")
	for {
		time.Sleep(3 * time.Second)
		worker.NewMariaWorker(mongo, maria).Restore()
	}
}

func SitenameWorker(mongo mongodbiface.DB, maria *sql.DB) {
	log.Println("Starting Sitename worker")
	for {
		time.Sleep(30 * time.Second)
		worker.NewSitenameWorker(mongo, maria).UpdateSitename()
	}
}

func NFSWorker(mongo mongodbiface.DB) {
	log.Println("Starting NFS worker")
	for {
		time.Sleep(3 * time.Second)
		worker.NewNFSWorker(mongo).Update()
	}
}

func CourseWorker(mongo mongodbiface.DB, maria *sql.DB) {
	log.Println("Starting Course worker")
	for {
		time.Sleep(5 * time.Second)
		worker.NewCourseWorker(mongo, maria).Update()
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
