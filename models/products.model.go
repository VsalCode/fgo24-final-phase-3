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

type Products struct {

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

func FindAllProducts() ([]Categories, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return []Categories{}, err
	}
	defer conn.Close()

	query := `
	SELECT p.id, p.code_product, p.name, p.image_url, p.purchase_price, p.selling_price, p.user_id 
	FROM products p 
	LEFT JOIN product_categories pc ON pc.id = p.category_id
	`
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