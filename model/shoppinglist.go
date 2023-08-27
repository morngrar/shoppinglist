package model

import "github.com/google/uuid"

type Item struct {
	Id   uint32 `bson:"uuid"`
	Name string `bson:"name"`
}

func NewItem(name string) Item {
	return Item{
		Id:   uuid.New().ID(),
		Name: name,
	}
}

type ShoppingList struct {
	Uuid  string `bson:"uuid"`
	Items []Item `bson:"items"`
}

func NewShoppingList() *ShoppingList {
	return &ShoppingList{Uuid: uuid.NewString(), Items: []Item{}}
}

func (sl ShoppingList) TemplateMapping() TemplateMapping {
	m := []ItemMapping{}
	for _, e := range sl.Items {
		m = append(m, ItemMapping{SLID: sl.Uuid, Id: e.Id, Name: e.Name})
	}
	return TemplateMapping{Items: m, Uuid: sl.Uuid}
}

type ItemMapping struct {
	SLID string
	Id   uint32
	Name string
}

type TemplateMapping struct {
	Items []ItemMapping
	Uuid  string
}

func NewItemMapping(sl *ShoppingList, item Item) ItemMapping {
	return ItemMapping{SLID: sl.Uuid, Id: item.Id, Name: item.Name}
}
