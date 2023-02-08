package moodle

import (
	"moodle/pkg/mongodbiface"
)

type MongoDB struct {
	collectionName string
	mongodbiface.MongoDB
}

func NewMongoDB(collectionName string, client mongodbiface.MongoDB) *MongoDB {
    return &MongoDB{
		collectionName: collectionName,
		MongoDB:        client,
	}
}
