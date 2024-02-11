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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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
	stmt, err := db.connection.Prepare(
		"select * from shoppinglists where uuid=$1",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare statement: %w", err)
	}

	sl := model.ShoppingList{}
	row := stmt.QueryRow(id)
	err = row.Scan(&sl.Id, &sl.Uuid)
	if err != nil {
		return nil, fmt.Errorf("Unable to scan row for getting shopping list: %w", err)
	}
	stmt.Close()

	stmt, err = db.connection.Prepare(
		"select * from items where shoppinglist_id=$1",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(sl.Id)
	if err != nil {
		return nil, fmt.Errorf("Unable to query items for shoppinglist %d: %w", sl.Id, err)
	}
	defer rows.Close()

	for rows.Next() {
		item := model.Item{}
		var meh int64
		err = rows.Scan(&item.Id, &item.Uuid, &meh, &item.Name, &item.Completed)
		if err != nil {
			log.Printf("Failed scanning row: %s", err)
			continue
		}

		sl.Items = append(sl.Items, item)
	}
	if err != nil {
		return nil, fmt.Errorf("failed scanning result set. Last error: %w", err)
	}

	return &sl, nil
}

// InsertShoppingList tries to insert a given list into the database, returns
// nil on success.
func (db PostgresHandle) InsertShoppingList(sl *model.ShoppingList) error {
	stmt, err := db.connection.Prepare("insert into shoppinglists (uuid) values ($1)")
	if err != nil {
		return fmt.Errorf("Unable to prepare statement: %w", err)
	}

	result, err := stmt.Exec(sl.Uuid)
	if err != nil {
		return fmt.Errorf("Unable to execute shopping list insert statement: %w", err)
	}
	stmt.Close()

	num, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unable to get rows affected: %w", err)
	}
	if num != 1 {
		return fmt.Errorf("Unexpected number of rows affected: %d", num)
	}

	return nil
}

// DeleteShoppingListByUuid attempts to delete a shopping list, returns any
// errors directly, nil on success.
func (db PostgresHandle) DeleteShoppingListByUuid(uuid string) error {
	return errors.New("not implemented")
}

// AddItemToShoppingList attempts to add an item to a shoppinglist of given
// uuid, returns any errors directly, nil on success.
func (db PostgresHandle) AddItemToShoppingList(uuid string, item model.Item) error {

	sl, err := db.GetShoppingListByID(uuid)
	if err != nil {
		return fmt.Errorf("Cannot retrieve shoppinglist for adding new item: %w", err)
	}

	stmt, err := db.connection.Prepare(
		"insert into items (uuid, shoppinglist_id, name, completed) values ($1, $2, $3, false)",
	)
	if err != nil {
		return fmt.Errorf("Unable to prepare statement: %w", err)
	}

	result, err := stmt.Exec(item.Uuid, sl.Id, item.Name)
	if err != nil {
		return fmt.Errorf("Unable to execute statement for adding item: %w", err)
	}
	stmt.Close()

	num, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unable to get rows affected: %w", err)
	}
	if num != 1 {
		return fmt.Errorf("Unexpected number of rows affected: %d", num)
	}

	return nil
}

// CompleteItemFromShoppingList sets the completed flag of a specified item
func (db PostgresHandle) CompleteItemFromShoppingList(slUuid string, itemUuid uint32) error {

	sl, err := db.GetShoppingListByID(slUuid)
	if err != nil {
		return fmt.Errorf("Cannot retrieve shoppinglist for completing item: %w", err)
	}

	stmt, err := db.connection.Prepare(
		"update items set completed = true where shoppinglist_id = $1 and uuid = $2",
	)
	if err != nil {
		return fmt.Errorf("Unable to prepare statement: %w", err)
	}

	result, err := stmt.Exec(sl.Id, itemUuid)
	if err != nil {
		return fmt.Errorf("Unable to execute statement for completing item: %w", err)
	}
	stmt.Close()

	num, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unable to get rows affected: %w", err)
	}
	if num != 1 {
		return fmt.Errorf("Unexpected number of rows affected: %d", num)
	}

	return nil
}

// RemoveItemFromShoppingList attempts to delete an item to a shoppinglist of
// given uuid, returns any errors directly, nil on success.
func (db PostgresHandle) RemoveItemFromShoppingList(slUuid string, itemUuid uint32) error {
	sl, err := db.GetShoppingListByID(slUuid)
	if err != nil {
		return fmt.Errorf("Cannot retrieve shoppinglist for deleting item: %w", err)
	}

	stmt, err := db.connection.Prepare(
		"delete from items where shoppinglist_id = $1 and uuid = $2",
	)
	if err != nil {
		return fmt.Errorf("Unable to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(sl.Id, itemUuid)
	if err != nil {
		return fmt.Errorf("Unable to execute statement for deleting item: %w", err)
	}

	num, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Unable to get rows affected: %w", err)
	}
	if num != 1 {
		return fmt.Errorf("Unexpected number of rows affected: %d", num)
	}

	return nil
}

func (db PostgresHandle) GetItemByID(slId string, itemId uint32) (*model.Item, error) {

	stmt, err := db.connection.Prepare(
		"select * from shoppinglists where uuid=$1",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare statement: %w", err)
	}

	sl := model.ShoppingList{}
	row := stmt.QueryRow(slId)
	err = row.Scan(&sl.Id, &sl.Uuid)
	if err != nil {
		return nil, fmt.Errorf("Unable to scan row for getting shopping list: %w", err)
	}
	stmt.Close()

	stmt, err = db.connection.Prepare(
		//"select (id, uuid, name, completed) from items where shoppinglist_id=$1 and uuid=$2",  // for some reason doesn't work
		"select * from items where shoppinglist_id=$1 and uuid=$2",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to prepare statement: %w", err)
	}
	defer stmt.Close()

	item := model.Item{}
	var meh int64
	row = stmt.QueryRow(sl.Id, itemId)
	err = row.Scan(&item.Id, &item.Uuid, &meh, &item.Name, &item.Completed)
	if err != nil {
		return nil, fmt.Errorf("Unable to scan row for getting item: %w", err)
	}

	return &item, nil
}
