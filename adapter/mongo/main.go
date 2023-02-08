package adapter

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAdapter interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}
