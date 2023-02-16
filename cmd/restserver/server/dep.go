package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"moodle/internal/dep"
	"moodle/internal/handler"
	"time"

	"moodle/config"
	moodleService "moodle/internal/core/service/moodle"
	Repo "moodle/internal/repository/moodle"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func initDependencies() *dep.Dep {
	d := &dep.Dep{}

	d.MongoRepository = Repo.NewMongoDB(NewMongoDB())
	d.MariaRepository = Repo.NewMariaDB(NewMariaDB())
	d.K8sRepository = Repo.NewK8sClient()
	d.MoodleService = moodleService.NewService(d.MongoRepository, d.MariaRepository, d.K8sRepository)
	d.MoodleHandler = handler.NewMoodleHandler(d.MoodleService)

	return d
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
	log.Println("Connected to MongoDB!")
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
	log.Println("Connected to MariaDB!")
	return db

}
