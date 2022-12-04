package worker

import (
	"context"
	"log"
	"moodle/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func checkMariaInstanceStatus(name string) (string, error) {
	// mongoclient
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// check instance status by name
	collection := client.Database("moodle").Collection("maria_instances")
	filter := bson.D{{Key: "name", Value: name}}
	var result bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result["status"].(string), nil
}


// k8s worker
type K8sWorker struct {
}

func NewK8sWorker() *K8sWorker {
	return &K8sWorker{}
}


// list maria id from maria_queue
func (k *K8sWorker) ListMariaQueue() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database("moodle").Collection("maria_queue")
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	var results []string
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result["id"].(string))
	}
	return results, nil
}


// get list of ACTIVE maria instances from maria id list
func (k *K8sWorker) GetActiveMariaInstances(mariaIds []string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGOURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	collection := client.Database("moodle").Collection("maria_instances")
	var results []string
	if len(mariaIds) > 0 {
		for _, mariaId := range mariaIds {
			filter := bson.M{"name": mariaId}
			var result bson.M
			err := collection.FindOne(context.Background(), filter).Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			if result["status"].(string) == "ACTIVE" {
				results = append(results, result["id"].(string))
			}
		}
	}

	return results, nil
}
