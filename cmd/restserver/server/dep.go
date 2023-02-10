package server

import (
	"context"
	"fmt"
	"log"
	"moodle/internal/dep"
	"moodle/internal/handler"

	"moodle/config"
	moodleService "moodle/internal/core/service/moodle"
	Repo "moodle/internal/repository/moodle"
	k8sRepo "moodle/pkg/client"

	"moodle/workers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func initDependencies() *dep.Dep {
	d := &dep.Dep{}

	d.MoodleRepository = Repo.NewMongoDB(NewMongoDB())
	d.MariaRepository = Repo.NewMariaDB(config.MARIAHOSTW, config.MARIAPORT, config.MARIAUSER, config.MARIAPASSWORD)
	d.K8sRepository = k8sRepo.NewK8sClient()
	d.MoodleService = moodleService.NewService(d.MoodleRepository, d.MariaRepository, d.K8sRepository)
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
	fmt.Println("Connected to MongoDB!")
	return client.Database(config.DBNAME)
}

func RunWorker() {
	workers.Start(NewMongoDB())
}
