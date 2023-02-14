package models

import (
	"context"
	"time"
)

type Order struct {
	ID            int       `json:"id"`
	ItemID        int       `json:"item_id"`
	TransactionID int       `json:"transaction_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CustomerID    int       `json:"customer_id"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func (m *DBModel) InsertOrder(ord Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	stmt := "INSERT INTO orders " +
		"( `item_id`, `transaction_id`, `status_id`, `quantity`, `amount`,`created_at`,`updated_at`,`customer_id`)" +
		" VALUES (?,?,?,?,?,?,?,?)"
	result, err := m.DB.ExecContext(ctx, stmt,
		ord.ItemID,
		ord.TransactionID,
		ord.StatusID,
		ord.Quantity,
		ord.Amount,
		time.Now(),
		time.Now(),
		ord.CustomerID,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
