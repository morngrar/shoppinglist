package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/morngrar/shoppinglist/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const databaseString = "shoppinglist"
const shoppinglistCollectionString = "shoppinglists"

// DatabaseHandle represents a living database connection
type DatabaseHandle struct {
	ctx    context.Context
	client *mongo.Client
	uri    string
}

// NewDatabaseHandle initiates a connection. Don't forget to call 'Disconnect'
// on it when done. Returns error on failure.
func NewDatabaseHandle() (DatabaseHandle, error) {
	dbh := DatabaseHandle{}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	url := os.Getenv("DB_HOST")

	protocol := "mongodb+srv"
	if url == "localhost" { // in case of testing locally, mongodb+srv wont work
		protocol = "mongodb"
	}

	if user == "" || password == "" || url == "" {
		return dbh, errors.New("environment not set up correctly")
	}

	uri := fmt.Sprintf(
		"%s://%s:%s@%s:27017/?authMechanism=DEFAULT",
		protocol,
		user,
		password,
		url,
	)
	log.Printf("Attempting to connect to mongo at %s", uri)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return dbh, fmt.Errorf("unable to connect to mongodb: %w", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return dbh, fmt.Errorf("mongo client ping failed: %w", err)
	}

	dbh.ctx = ctx
	dbh.client = client
	dbh.uri = uri
	return dbh, nil
}

// Disconnect ends the connection of the handle.
func (db DatabaseHandle) Disconnect() {
	db.client.Disconnect(db.ctx)
}

// GetShoppingListByID takes a shopping list UUID and tries to retrieve that it
// from the database.
func (db DatabaseHandle) GetShoppingListByID(id string) (*model.ShoppingList, error) {
	var sl model.ShoppingList
	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)
	err := collection.FindOne(db.ctx, bson.M{"uuid": id}).Decode(&sl)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve a shopping list with id <%s>: %w",
			id, err,
		)
	}

	return &sl, nil
}

// InsertShoppingList tries to insert a given list into the database, returns
// nil on success.
func (db DatabaseHandle) InsertShoppingList(sl *model.ShoppingList) error {

	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)

	_, err := collection.InsertOne(db.ctx, sl)
	if err != nil {
		return fmt.Errorf(
			"failed to do shopping list insertion: %w", err,
		)
	}

	return nil
}

// DeleteShoppingListByUuid attempts to delete a shopping list, returns any
// errors directly, nil on success.
func (db DatabaseHandle) DeleteShoppingListByUuid(uuid string) error {

	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)

	_, err := collection.DeleteOne(db.ctx, bson.M{"uuid": uuid})

	return err
}

// AddItemToShoppingList attempts to add an item to a shoppinglist of given
// uuid, returns any errors directly, nil on success.
//
// Ref: https://www.mongodb.com/docs/manual/reference/operator/update/push/
func (db DatabaseHandle) AddItemToShoppingList(uuid string, item model.Item) error {

	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)

	_, err := collection.UpdateOne(
		db.ctx,
		bson.M{"uuid": uuid},
		bson.M{"$push": bson.M{"items": item}},
	)

	return err
}

// CompleteItemFromShoppingList sets the completed flag of a specified item
func (db DatabaseHandle) CompleteItemFromShoppingList(slUuid string, itemUuid uint32) error {
	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)

	update := bson.M{
		"$set": bson.M{
			"items.$.completed": true, // Set completed to true (or false if you want to mark it as incomplete)
		},
	}

	data, err := collection.UpdateOne(
		db.ctx,
		bson.M{"uuid": slUuid, "items.uuid": itemUuid},
		update,
	)

	fmt.Println(data)

	return err
}

// RemoveItemFromShoppingList attempts to delete an item to a shoppinglist of
// given uuid, returns any errors directly, nil on success.
//
// Ref: https://www.mongodb.com/docs/manual/reference/operator/update/pull/
func (db DatabaseHandle) RemoveItemFromShoppingList(slUuid string, itemUuid uint32) error {

	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)

	_, err := collection.UpdateOne(
		db.ctx,
		bson.M{"uuid": slUuid},
		bson.M{"$pull": bson.M{"items": bson.M{"uuid": itemUuid}}},
	)

	return err
}

func (db DatabaseHandle) GetItemByID(slId string, itemId uint32) (*model.Item, error) {

	collection := db.client.
		Database(databaseString).
		Collection(shoppinglistCollectionString)

	var sl model.ShoppingList

	projection := bson.M{
		"items": bson.M{
			"$elemMatch": bson.M{
				"uuid": itemId,
			},
		},
	}

	err := collection.FindOne(
		db.ctx,
		bson.M{"uuid": slId},
		options.FindOne().SetProjection(projection),
	).Decode(&sl)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to retrieve a shopping list item with id <%d> from <%s>: %w",
			itemId, slId, err,
		)
	}

	item := sl.Items[0]

	return &item, nil
}
