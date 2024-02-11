package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/morngrar/shoppinglist/model"

	_ "github.com/lib/pq"
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

const (
	defaultHost     = "localhost"
	defaultPort     = "5432"
	defaultUser     = "postgres"
	defaultPassword = "postgres"
	defaultDatabase = "shoppinglist"
)

// PostgresHandle represents a living database connection
type PostgresHandle struct {
	connection *sql.DB
}

func NewPostgresHandle() (PostgresHandle, error) {
	var err error
	handle := PostgresHandle{}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = defaultHost
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = defaultPort
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = defaultUser
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = defaultPassword
	}

	database := os.Getenv("DB_DATABASE")
	if database == "" {
		database = defaultDatabase
	}

	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmod=disable",
		host, port, user, password, database,
	)

	handle.connection, err = sql.Open("postgres", connString)
	if err != nil {
		log.Printf("Failed to connect to postgres instance using: %q", connString)
		return handle, err
	}

	err = handle.connection.Ping()
	if err != nil {
		log.Printf(
			"Failed to ping postgres instance after connecting with: %q",
			connString,
		)
		handle.connection.Close()
		return handle, err
	}

	return handle, nil
}

// Disconnect ends the connection of the handle.
func (db PostgresHandle) Disconnect() {
	db.connection.Close()
}

// GetShoppingListByID takes a shopping list UUID and tries to retrieve that it
// from the database.
func (db PostgresHandle) GetShoppingListByID(id string) (*model.ShoppingList, error) {
	return nil, errors.New("Not implemented")
}

// InsertShoppingList tries to insert a given list into the database, returns
// nil on success.
func (db PostgresHandle) InsertShoppingList(sl *model.ShoppingList) error {
	return errors.New("Not implemented")
}

// DeleteShoppingListByUuid attempts to delete a shopping list, returns any
// errors directly, nil on success.
func (db PostgresHandle) DeleteShoppingListByUuid(uuid string) error {
	return errors.New("not implemented")
}

// AddItemToShoppingList attempts to add an item to a shoppinglist of given
// uuid, returns any errors directly, nil on success.
func (db PostgresHandle) AddItemToShoppingList(uuid string, item model.Item) error {
	return errors.New("not implemented")
}

// CompleteItemFromShoppingList sets the completed flag of a specified item
func (db PostgresHandle) CompleteItemFromShoppingList(slUuid string, itemUuid uint32) error {
	return errors.New("not implemented")
}

// RemoveItemFromShoppingList attempts to delete an item to a shoppinglist of
// given uuid, returns any errors directly, nil on success.
func (db PostgresHandle) RemoveItemFromShoppingList(slUuid string, itemUuid uint32) error {
	return errors.New("not implemented")
}

func (db PostgresHandle) GetItemByID(slId string, itemId uint32) (*model.Item, error) {
	return nil, errors.New("not implemented")
}
