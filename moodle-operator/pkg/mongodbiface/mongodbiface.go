package mongodbiface

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

