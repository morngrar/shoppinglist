package model

import "github.com/google/uuid"

type Item struct {
	Id        uint32 `bson:"uuid"`
	Name      string `bson:"name"`
	Completed bool   `bson:"completed"`
}

func NewItem(name string) Item {
	return Item{
		Id:        uuid.New().ID(),
		Name:      name,
		Completed: false,
	}
}

func (i Item) TemplateMapping(slId string) ItemMapping {
	return ItemMapping{
		SLID:      slId,
		Id:        i.Id,
		Name:      i.Name,
		Completed: i.Completed,
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
		m = append(
			m,
			e.TemplateMapping(sl.Uuid),
		)
	}
	return TemplateMapping{Items: m, Uuid: sl.Uuid}
}

type ItemMapping struct {
	SLID      string
	Id        uint32
	Name      string
	Completed bool
}

type TemplateMapping struct {
	Items []ItemMapping
	Uuid  string
}

func NewItemMapping(sl *ShoppingList, item Item) ItemMapping {
	return item.TemplateMapping(sl.Uuid)
}
