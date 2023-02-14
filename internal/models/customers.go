package models

import (
	"context"
	"time"
)

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"_"`
	UpdatedAt time.Time `json:"_"`
}

func (m *DBModel) InsertCustomer(cst Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	stmt := "INSERT INTO customers " +
		"(`first_name`, `last_name`, `email`,`created_at`,`updated_at`)" +
		" VALUES (?,?,?,?,?)"
	result, err := m.DB.ExecContext(ctx, stmt,
		cst.FirstName,
		cst.LastName,
		cst.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
