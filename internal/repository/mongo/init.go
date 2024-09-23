package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MClient struct {
	client *mongo.Client
	dbname string
}

func NewMongoClient(url, dbname string) (*MClient, error) {

	var cl MClient

	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(url))
	if err != nil {
		return nil, err
	}

	cl.client = client
	cl.dbname = dbname

	return &cl, nil
}
