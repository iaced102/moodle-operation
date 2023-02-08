package mongodbiface

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB interface {
	Connect() (*mongo.Client, error)
	Ping() error
	Database() *mongo.Database
	Collection() *mongo.Collection
}

