package model

import "github.com/google/uuid"

// Item is the database model of a shopping list item
type Item struct {
	Id        int64
	Uuid      uint32
	Name      string
	Completed bool
}

func NewItem(name string) Item {
	return Item{
		Uuid:      uuid.New().ID(),
		Name:      name,
		Completed: false,
	}
}

func (i Item) TemplateMapping(slId string) ItemMapping {
	return ItemMapping{
		SLID:      slId,
		Id:        i.Uuid,
		Name:      i.Name,
		Completed: i.Completed,
	}
}

// ShoppingList is the database model of a shoppingn list
type ShoppingList struct {
	Id    int64
	Uuid  string
	Items []Item
}

func NewShoppingList() *ShoppingList {
	return &ShoppingList{Uuid: uuid.NewString(), Items: []Item{}}
}

func (sl ShoppingList) TemplateMapping() TemplateMapping {
	m := []ItemMapping{}
	for _, e := range sl.Items {
		m = append(
			m,
			e.TemplateMapping(sl.Uuid),
		)
	}
	return TemplateMapping{Items: m, Uuid: sl.Uuid}
}

// ItemMapping is the mapping of a shopping list item to an object usable in
// the template engine
type ItemMapping struct {
	SLID      string
	Id        uint32
	Name      string
	Completed bool
}

// TemplateMapping is the mapping of a shopping list to an object usable in
// the template engine
type TemplateMapping struct {
	Items []ItemMapping
	Uuid  string
}

func NewItemMapping(sl *ShoppingList, item Item) ItemMapping {
	return item.TemplateMapping(sl.Uuid)
}
