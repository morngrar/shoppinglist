package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/morngrar/shoppinglist/database"
	"github.com/morngrar/shoppinglist/model"
)

var latestId int64 = 0

func getId() int64 {
	latestId++
	return latestId - 1
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDatabaseHandle()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Disconnect()

	sl := model.NewShoppingList()
	err = db.InsertShoppingList(sl)
	if err != nil {
		log.Printf("Failed to insert new shopping list: %s", err)
		http.Error(w, "Failed to create shopping list", http.StatusInternalServerError)
		return
	}

	newPage := fmt.Sprintf("/%s", sl.Uuid)
	w.Header().Add("Location", newPage)
	w.WriteHeader(http.StatusFound)
}

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// define handlers
	r.Get("/", indexHandler)
	r.Get("/{shoppingListId}", shoppingListGetHandler)
	r.Post("/{shoppingListId}/add-item/", shoppingListAddHandler)
	r.Get("/{shoppingListId}/complete-item/{id}", itemDoneHandler)
	r.Delete("/{shoppingListId}/complete-item/{id}", itemDeleteHandler)

	port := "8080"
	fmt.Println("Running on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
