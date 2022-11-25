package worker

import (
	"context"
	"log"
	lb "moodle/client/lb"
	"moodle/handler"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var client *lb.LBClient

type LBWoker struct {
}


func  NewLBWorker() *LBWoker {
	return &LBWoker{}
}

// list lb then insert to mongodb
func (m *LBWoker) GetLBInstances() {
	client = lb.NewLBClient()
	// get lb instances
	lbInstances := client.GetLBInstances()
	// connect to mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(handler.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// get collection
	collection := client.Database(handler.DBNAME).Collection("lb")
	// insert to mongodb
	for _, lbInstance := range lbInstances {
		_, err := collection.InsertOne(ctx, lbInstance)
		if err != nil {
			log.Fatal(err)
		}
	}
}
