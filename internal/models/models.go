package models

import (
	"context"
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

type Order struct {
	ID            int       `json:"id"`
	ProductID     int       `json:"item_id"`
	TransactionID int       `json:"transaction_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"_"`
	UpdatedAt     time.Time `json:"_"`
}

type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

type Transactions struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_Four"`
	BankReturnCode      string    `json:"bank_Return_Code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	CreatedAt           time.Time `json:"_"`
	UpdatedAt           time.Time `json:"_"`
}
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

// Item is the type for all Items
type Item struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"_"`
	UpdatedAt      time.Time `json:"_"`
}

func (m *DBModel) GetItem(id int) (Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var product Item
	row := m.DB.QueryRowContext(ctx,
		"select id,name,description,inventory_level,price,"+
			"coalesce (image,''),created_at,updated_at from items where id =?", id)
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.InventoryLevel,
		&product.Price,
		&product.Image,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return product, err
	}
	return product, nil
}
func (m *DBModel) GetAllItems() ([]Item, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var Items []Item
	rows, err := m.DB.QueryContext(ctx,
		"select id,name,description,inventory_level,price,"+
			"coalesce (image,''),created_at,updated_at from items ")
	for rows.Next() {
		var product Item
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.InventoryLevel,
			&product.Price,
			&product.Image,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return Items, err
		}
		Items = append(Items, product)
	}
	if err = rows.Err(); err != nil {
		return Items, err
	}

	return Items, nil
}
