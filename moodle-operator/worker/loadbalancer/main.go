package worker

import (
	"context"
	"log"
	mongoadapter "moodle/adapter/mongo"
	lb "moodle/client/lb"
	"moodle/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type LBWoker struct {
}


func  NewLBWorker() *LBWoker {
	return &LBWoker{}
}

// new mongoadapter
func (l *LBWoker) NewMongoAdapter() *mongoadapter.Adapter {
		ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.NewClient(
		options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatalf("Error creating mongo client: %+v", err)
	}
	defer client.Disconnect(ctx)
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %+v", err)
	}
	return mongoadapter.New(client.Database(config.DBNAME))
}

// list lb then insert to mongodb
func (m *LBWoker) GetLBInstances() error {
	var client *lb.LBClient = lb.NewLBClient()

	// list instances
	instances, err := client.GetLBInstances()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
	   log.Fatal(err)
	}
	err = mongoClient.Connect(ctx)
	if err != nil {
	   log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)
	collection := mongoClient.Database(config.DBNAME).Collection("lb_instances")





	// for lb in collection if not in instances then delete
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// if not in instances then delete
		if !inInstances(result["name"].(string), instances) {
			_, err := collection.DeleteOne(ctx, bson.M{"name": result["name"]})
			if err != nil {
				log.Fatal(err)
			}
		}
	}


	for _, instance := range instances {
		filter := bson.M{"name": instance.Name}
		update := bson.M{"$set": instance}
		// insert to mongodb, update if exists
		_, err = collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal(err)
		}
	}
	
	return nil
}


func inInstances(name string, instances []*lb.LoadBalancer) bool {
	for _, instance := range instances {
		if instance.Name == name {
			return true
		}
	}
	return false
}

