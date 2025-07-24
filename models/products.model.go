package models

import (
	"context"
	"fmt"
	"nashta_inventory/db"
	"nashta_inventory/dto"

	"github.com/jackc/pgx/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Categories struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Products struct {
	Id            int     `json:"id"`
	CodeProduct   string  `json:"codeProduct" db:"code_product"`
	Name          string  `json:"name"`
	ImageURL      string  `json:"imageUrl" db:"image_url"`
	PurchasePrice float64 `json:"purchasePrice" db:"purchase_price"`
	SellingPrice  float64 `json:"sellingPrice" db:"selling_price"`
	Quantity      int     `json:"quantity"`
	User          string  `json:"user"`
	Category      string  `json:"categories"`
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

func CreateNewProducts(req dto.ProductRequest) (Products, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return Products{}, err
	}
	defer conn.Close()

	id, err := gonanoid.New(6)
	if err != nil {
		return Products{}, err
	}
	code := fmt.Sprintf("PRDK-%s", id)

	query := `
	INSERT INTO products (code_product, name, image_url, purchase_price, selling_price, quantity, user_id, category_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`

	var productId int
	err = conn.QueryRow(
		context.Background(), 
		query, 
		code, req.Name, req.ImageURL, req.PurchasePrice, req.SellingPrice, req.Quantity, 
		req.UserId, req.CategoryId).Scan(&productId)
	if err != nil {
		return Products{}, err
	}

	querySelect := `
	SELECT p.id, p.code_product, p.name, p.image_url, p.purchase_price, p.selling_price, p.quantity, u.name as user, pc.name as category
	FROM products p
	JOIN users u ON u.id = p.user_id
	JOIN product_categories pc ON pc.id = p.category_id  
	WHERE p.id = $1
	`
	row, err := conn.Query(context.Background(), querySelect, productId)
	if err != nil {
		return Products{}, err
	}

	result, err := pgx.CollectOneRow[Products](row, pgx.RowToStructByName)
	if err != nil {
		return Products{}, err
	}

	return result, nil
}