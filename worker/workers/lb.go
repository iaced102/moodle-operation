package worker

import (
	"context"
	"log"
	adapter "moodle/adapter/mongo"
	lb "moodle/pkg/client"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type LBWorker struct {
    adapter adapter.MongoAdapter
}

func NewLBWorker(m adapter.MongoAdapter) *LBWorker{
    return &LBWorker{
		adapter: m,
    }
}

// list lb then insert to mongodb
func (m *LBWorker) GetLBInstances() error {
	var client *lb.LBClient = lb.NewLBClient()
	// get token from tokens collection
	var result bson.M
	ctx := context.Background()
	tokenCollection := m.adapter.Collection("tokens")
	err := tokenCollection.FindOne(context.Background(), bson.M{}).Decode(&result)
	token := result["token"].(string)

	// list instances
	instances, err := client.GetLBInstances(token)
	if err != nil {
		log.Fatal(err)
	}
	collection := m.adapter.Collection("lb_instances")

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

