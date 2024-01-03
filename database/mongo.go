package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var (
	mongoClient *mongo.Client
	ctx         = context.Background()
	dbName      string
)

func StartMongoDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("you must set your 'MONGODB_URI' environmental variable. See\\n\\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		return errors.New("you must set your 'DATABASE' environmental variable")
	} else {
		dbName = database
	}

	var err error
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}
	return nil
}

func CloseMongoDB() {
	err := mongoClient.Disconnect(ctx)
	if err != nil {
		panic(err)
	}
}

func GetBucket() *gridfs.Bucket {
	bucket, err := gridfs.NewBucket(mongoClient.Database(dbName))
	if err != nil {
		panic(err)
	}

	return bucket
}

func GetCollection(col string) *mongo.Collection {
	return mongoClient.Database(dbName).Collection(col)
}

func GetMongoClient() *mongo.Client {
	return mongoClient
}