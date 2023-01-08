package worker

import (
	"context"
	"log"
	adapter "moodle/adapter/mongo"
	maria "moodle/client/maria"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type MariaWorker struct {
    adapter adapter.MongoAdapter
}

func NewMariaWorker(m adapter.MongoAdapter) *MariaWorker{
    return &MariaWorker{
		adapter: m,
    }
}

// list maria instances then insert to mongodb
func (m *MariaWorker) GetMariaInstances() error {

	var client *maria.MariaClient = maria.NewMariaClient()

	// list instances
	instances, err := client.GetMariaInstances()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	collection := m.adapter.Collection("maria_instances")

	// for maria instance in collection if not in instances then delete
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


func inInstances(name string, instances []*maria.CloudDatabaseInstance ) bool {
	for _, instance := range instances {
		if instance.Name == name {
			return true
		}
	}
	return false
}
