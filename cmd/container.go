package cmd

import (
	"context"
	"crypto/tls"

	"github.com/d-leme/tradew-users/pkg/core"
	"github.com/d-leme/tradew-users/pkg/users"
	"github.com/d-leme/tradew-users/pkg/users/mongodb"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Container contains all depencies from our api
type Container struct {
	Settings *core.Settings

	MongoClient *mongo.Client

	UserRepository users.Repository
	UserService    users.Service
	UserController users.Controller
}

// NewContainer creates new instace of Container
func NewContainer(settings *core.Settings) *Container {

	container := new(Container)

	container.Settings = settings

	container.MongoClient = connectMongoDB(settings.MongoDB)

	container.UserRepository = mongodb.NewRepository(container.MongoClient, settings.MongoDB.Database)
	container.UserService = users.NewService(settings, container.UserRepository)
	container.UserController = users.NewController(container.UserService)

	return container
}

// Controllers maps all routes and exposes them
func (c *Container) Controllers() []core.Controller {
	return []core.Controller{
		&c.UserController,
	}
}

// Close terminates every opened resources
func (c *Container) Close() {
	c.MongoClient.Disconnect(context.Background())
}

func connectMongoDB(conf *core.MongoDBConfig) *mongo.Client {
	ctx := context.Background()

	client, err := mongo.NewClient(options.
		Client().
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		ApplyURI(conf.ConnectionString))

	if err != nil {
		logrus.
			WithError(err).
			Fatal("error creating MongoDB client")
	}

	if err := client.Connect(ctx); err != nil {
		logrus.
			WithError(err).
			Fatal("error connecting to MongoDB")
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.
			WithError(err).
			Fatal("error pinging MongoDB")
	}

	return client
}
