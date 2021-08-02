package mongodb

import (
	"context"
	"time"

	"github.com/d-leme/tradew-users/pkg/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"gopkg.in/mgo.v2/bson"
)

type repositoryMongoDB struct {
	collection *mongo.Collection
}

// NewRepository ...
func NewRepository(client *mongo.Client, database string) users.Repository {
	repository := &repositoryMongoDB{client.Database(database).Collection("users")}
	repository.createIndex()

	return repository
}

// Insert ...
func (repository *repositoryMongoDB) Insert(ctx context.Context, p *users.User) error {
	if _, err := repository.collection.InsertOne(ctx, p); err != nil {
		return err
	}

	return nil
}

// GetByEmail ...
func (repository *repositoryMongoDB) GetByEmail(ctx context.Context, email string) (*users.User, error) {
	var r *users.User

	err := repository.collection.FindOne(ctx, bson.M{"email": email}).Decode(&r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (repository *repositoryMongoDB) createIndex() {
	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	filterEmail := mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "email", Value: bsonx.Int32(-1)},
		},
		Options: options.Index().SetUnique(false),
	}

	repository.collection.Indexes().CreateOne(ctx, filterEmail)
}
