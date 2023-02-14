package models

import (
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}
type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
type ApiResponse struct {
	OK  bool   `json:"ok"`
	ID  int    `json:"id"`
	MSG string `json:"message"`
}
