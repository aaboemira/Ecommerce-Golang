package models

import "time"

type DBModel struct {
}
type models struct {
	db DBModel
}

type Item struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory"`
	Price          int       `json:"price"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}
