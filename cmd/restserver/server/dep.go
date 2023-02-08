package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"moodle/config"
	"moodle/internal/dep"

	moodleService "moodle/internal/core/service/moodle"
	"moodle/internal/handler"
	moodleRepo "moodle/internal/repository/moodle"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionName = "moodles"
	mariadbName       = ""
)

func initDependencies() *dep.Dep {
	d := &dep.Dep{}

	d.MongoDB = NewMongoDB()
	d.MariaDB = NewMariaDB()

	d.MoodleRepository = moodleRepo.NewMongoDB(CollectionName, d.MongoDB)
	d.MariaRepository = moodleRepo.NewMariaDB(mariadbName, d.MariaDB)
	d.MoodleService = moodleService.NewService(d.MoodleRepository, d.MariaRepository)
	d.MoodleHandler = handler.NewGameHandler(d.MoodleService, d.MariaService)

	return d
}


func NewMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}
	return client
}


func NewMariaDB() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)", config.MARIAUSER, config.MARIAPASSWORD, config.MARIAHOSTW, config.MARIAPORT))
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	return db
}
