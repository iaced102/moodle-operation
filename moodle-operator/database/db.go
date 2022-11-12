package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type DBClient struct {
	db DB
}

func NewDB(d DB) *DBClient {
	return &DBClient{
		db: d,
	}
}

// findone
func (d *DBClient) FindOne(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return d.db.Collection(collection).FindOne(ctx, filter, opts...)
}

// find
func (d *DBClient) Find(ctx context.Context, collection string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return d.db.Collection(collection).Find(ctx, filter, opts...)
}

// insertone
func (d *DBClient) InsertOne(ctx context.Context, collection string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return d.db.Collection(collection).InsertOne(ctx, document, opts...)
}

// insertmany
func (d *DBClient) InsertMany(ctx context.Context, collection string, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return d.db.Collection(collection).InsertMany(ctx, documents, opts...)
}


// updateone
func (d *DBClient) UpdateOne(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return d.db.Collection(collection).UpdateOne(ctx, filter, update, opts...)
}

// updatemany
func (d *DBClient) UpdateMany(ctx context.Context, collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return d.db.Collection(collection).UpdateMany(ctx, filter, update, opts...)
}

// deleteone
func (d *DBClient) DeleteOne(ctx context.Context, collection string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return d.db.Collection(collection).DeleteOne(ctx, filter, opts...)
}

// deletemany
func (d *DBClient) DeleteMany(ctx context.Context, collection string, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return d.db.Collection(collection).DeleteMany(ctx, filter, opts...)
}
