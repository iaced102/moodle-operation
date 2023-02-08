package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"moodle/config"
	"moodle/internal/dep"

	gameService "github.com/matiasvarela/minesweeper-API/internal/core/service/game"
	"github.com/matiasvarela/minesweeper-API/internal/handler"
	gameRepo "github.com/matiasvarela/minesweeper-API/internal/repository/game"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbCollection = "moodles"
)

func initDependencies() *dep.Dep {
	d := &dep.Dep{}

	d.MongoDB = NewMongoDB()
	d.MariaDB = NewMariaDB()

	d.MoodleRepository = gameRepo.NewDynamoDB(mongodbCollection, d.MongoDB)
	d.GameService = gameService.NewService(rnd, clk, d.GameRepository)
	d.GameHandler = handler.NewGameHandler(d.GameService)

	return d
}


func NewMongoDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database(config.DBNAME)
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
