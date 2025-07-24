package models

import (
	"context"
	"nashta_inventory/db"

	"github.com/jackc/pgx/v5"
)

type Categories struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

func FindAllCategories() ([]Categories, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return []Categories{}, err
	}
	defer conn.Close()

	query := `SELECT id, name, description FROM product_categories`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return []Categories{}, err
	}

	categories, err := pgx.CollectRows[Categories](rows, pgx.RowToStructByName)
	if err != nil {
		return []Categories{}, err
	}

	return categories, nil
}