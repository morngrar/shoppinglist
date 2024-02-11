package database

import (
	"context"
	"fmt"

	"github.com/morngrar/shoppinglist/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	AddItemToShoppingList(uuid string, item model.Item) error
	CompleteItemFromShoppingList(slUuid string, itemUuid uint32) error
	DeleteShoppingListByUuid(uuid string) error
	Disconnect()
	GetItemByID(slId string, itemId uint32) (*model.Item, error)
	GetShoppingListByID(uuid string) (*model.ShoppingList, error)
	InsertShoppingList(sl *model.ShoppingList) error
	RemoveItemFromShoppingList(slUuid string, itemUuid uint32) error
}

// PostgresHandle represents a living database connection
type PostgresHandle struct {
	ctx    context.Context
	client *mongo.Client
	uri    string
}

// Disconnect ends the connection of the handle.
func (db PostgresHandle) Disconnect() {
	db.client.Disconnect(db.ctx)
}

// GetShoppingListByID takes a shopping list UUID and tries to retrieve that it
// from the database.
func (db PostgresHandle) GetShoppingListByID(id string) (*model.ShoppingList, error) {
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
func (db PostgresHandle) InsertShoppingList(sl *model.ShoppingList) error {

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
func (db PostgresHandle) DeleteShoppingListByUuid(uuid string) error {

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
func (db PostgresHandle) AddItemToShoppingList(uuid string, item model.Item) error {

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
func (db PostgresHandle) CompleteItemFromShoppingList(slUuid string, itemUuid uint32) error {
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
func (db PostgresHandle) RemoveItemFromShoppingList(slUuid string, itemUuid uint32) error {

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

func (db PostgresHandle) GetItemByID(slId string, itemId uint32) (*model.Item, error) {

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
