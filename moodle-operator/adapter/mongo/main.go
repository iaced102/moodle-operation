package adapter

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoAdapter interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type Adapter struct {
	adapter MongoAdapter
}

func New(m MongoAdapter) *Adapter {
	return &Adapter{
		adapter: m,
	}
}
