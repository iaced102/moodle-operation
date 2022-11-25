package worker

import (
	"context"
	maria "moodle/client/maria"
	"moodle/handler"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



type MariaWorker struct {
}



func NewMariaWorker() *MariaWorker {
	return &MariaWorker{}
}


// list maria instances then insert to mongodb
func (m *MariaWorker) GetMariaInstances() error {

	var client *maria.MariaClient = maria.NewMariaClient()

	// list instances
	instances, err := client.ListInstances()
	if err != nil {
		return err
	}


	// insert to mongodb, update if exists
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(handler.MONGOURI))
	if err != nil {
		return err
	}
	err = mongoClient.Connect(context.Background())
	if err != nil {
		return err
	}
	collection := mongoClient.Database(handler.DBNAME).Collection("maria_instances")
	for _, instance := range instances {
		filter := bson.M{"name": instance.Name}
		update := bson.M{"$set": instance}
		_, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}



	return nil


}


