package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/companieshouse/innovation-accounts-statistics-tool/config"
	"github.com/companieshouse/innovation-accounts-statistics-tool/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionClient provides an interface by which to interact with a database.
type TransactionClient interface {
	GetAccountsTransactions(dataDescription string) (*[]models.Transaction, error)
	Shutdown()
}

// TransactionDatabaseClient is a concrete implementation of the Client interface.
type TransactionDatabaseClient struct {
	db MongoDatabaseInterface
}

// NewTransactionDatabaseClient returns a new implementation of the Client interface.
func NewTransactionDatabaseClient(cfg *config.Config) TransactionClient {
	return &TransactionDatabaseClient{
		db: getMongoDatabase(cfg.TransactionsMongoDBURL, cfg.TransactionsMongoDBDatabase),
	}
}

var mgoClient *mongo.Client

func getMongoClient(mongoDBURL string) *mongo.Client {

	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	client, err := mongo.Connect(ctx, clientOptions)

	// the program must bail out here if failing to establish a connection to the db, as this will run on application start-up
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// cache mongo client here, in preparation for disconnect on application shutdown
	mgoClient = client

	// check we can connect to the mongodb instance - again, bail out on failure
	pingContext, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()
	err = client.Ping(pingContext, nil)
	if err != nil {
		log.Error("ping to mongodb timed out. please check the connection to mongodb and that it is running")
		os.Exit(1)
	}

	log.Info("connected to mongodb successfully")

	return client
}

func getMongoDatabase(mongoDBURL, databaseName string) MongoDatabaseInterface {
	return getMongoClient(mongoDBURL).Database(databaseName)
}

// MongoDatabaseInterface is an interface that describes the mongodb driver.
type MongoDatabaseInterface interface {
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

// GetAccountsTransactions returns a slice of Transaction which are retrieved from a mongoDB.
func (c *TransactionDatabaseClient) GetAccountsTransactions(dataDescription string) (*[]models.Transaction, error) {

	entities := make([]models.Transaction, 0)

	collection := c.db.Collection("transactions")
	cur, err := collection.Find(context.Background(), bson.M{"data.status": "closed", "data.description": dataDescription})

	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {

		var entity models.Transaction
		err = cur.Decode(&entity)

		if err != nil {
			return nil, err
		}

		entities = append(entities, entity)
	}

	return &entities, nil
}

// Shutdown is a hook that can be used to clean up db resources.
func (c *TransactionDatabaseClient) Shutdown() {
	log.Info("Attempting to close the db connection thread pool")
	if mgoClient != nil {
		err := mgoClient.Disconnect(context.Background())
		if err != nil {
			log.Error(fmt.Sprintf("Failed to disconnect from the mongodb: %s", err))
			return
		}
		log.Info("disconnected from mongodb successfully")
	}
}
