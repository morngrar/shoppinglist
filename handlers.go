package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/morngrar/shoppinglist/database"
	"github.com/morngrar/shoppinglist/model"
)

const shoppingListIdString = "shoppingListId"

func shoppingListGetHandler(w http.ResponseWriter, r *http.Request) {

	shoppingListId := chi.URLParam(r, shoppingListIdString)

	tmpl, err := template.ParseFiles("templates/shopping-list.html")
	if err != nil {
		log.Printf("Failed to parse shopping-list.html: %s", err)
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	db, err := database.NewDatabaseHandle()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Disconnect()

	sl, err := db.GetShoppingListByID(shoppingListId)
	if err != nil {
		log.Printf("Failed to get shopping list from db: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, sl.TemplateMapping())
	if err != nil {
		log.Printf("Failed to execute template: %s", err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func shoppingListAddHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	log.Printf("name: %s", name)
	tmpl, err := template.ParseFiles("templates/shopping-list.html")
	if err != nil {
		log.Printf("Failed to parse shopping-list.html: %s", err)
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	shoppingListId := chi.URLParam(r, shoppingListIdString)

	db, err := database.NewDatabaseHandle()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Disconnect()

	sl, err := db.GetShoppingListByID(shoppingListId)
	if err != nil {
		log.Printf("Failed to get shopping list from db: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	newItem := model.NewItem(name)
	sl.Items = append(sl.Items, newItem)

	err = db.AddItemToShoppingList(shoppingListId, newItem)
	if err != nil {
		log.Printf("Failed adding item: %s", err)
		http.Error(w, "Unable to add item", http.StatusInternalServerError)
		return
	}

	w.Header().Add("HX-Trigger", "refresh-form")

	err = tmpl.ExecuteTemplate(w, "item-list-element", model.NewItemMapping(sl, newItem))
	if err != nil {
		log.Printf("Failed to execute template: %s", err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

func itemCompleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		log.Printf("Did NOT get DELETE on `complete-item` endpoint: %v", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Printf("Received invalid id on `complete-item` endpoint: %s", err)
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	log.Printf("completing id: %d", id)

	shoppingListId := chi.URLParam(r, shoppingListIdString)

	db, err := database.NewDatabaseHandle()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Disconnect()

	_, err = db.GetShoppingListByID(shoppingListId)
	if err != nil {
		log.Printf("Failed to get shopping list from db: %s", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = db.RemoveItemFromShoppingList(shoppingListId, uint32(id))
	if err != nil {
		log.Printf("Failed to item from db: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
