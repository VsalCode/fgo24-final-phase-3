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
	User          string  `json:"user,omitempty" db:"user"` 
	Category      string  `json:"categories,omitempty" db:"category"`
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

func FindAllProducts() ([]Products, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return []Products{}, err
	}
	defer conn.Close()

	query := `
	SELECT 
		p.id, p.code_product, p.name, p.image_url, 
		p.purchase_price, p.selling_price, p.quantity, 
		u.name AS user, 
		pc.name AS category 
	FROM products p 
	LEFT JOIN product_categories pc ON pc.id = p.category_id
	LEFT JOIN users u ON u.id = p.user_id 
	`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return []Products{}, err
	}

	products, err := pgx.CollectRows[Products](rows, pgx.RowToStructByName) // Changed categories to products for clarity
	if err != nil {
		return []Products{}, err
	}

	return products, nil
}

func CreateNewProducts(req dto.ProductRequest, userId int) (Products, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return Products{}, err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return Products{}, err
	}
	defer tx.Rollback(context.Background()) 

	id, err := gonanoid.New(6)
	if err != nil {
		return Products{}, err
	}
	code := fmt.Sprintf("BRG-%s", id)

	insertQuery := `
	INSERT INTO products (
		code_product, name, image_url, 
		purchase_price, selling_price, 
		quantity, user_id, category_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`

	var productId int
	err = tx.QueryRow(
		context.Background(),
		insertQuery,
		code, 
		req.Name, 
		req.ImageURL, 
		req.PurchasePrice, 
		req.SellingPrice, 
		req.Quantity, 
		userId, 
		req.CategoryId,
	).Scan(&productId)
	
	if err != nil {
		return Products{}, err
	}

	trxQuery := `
	INSERT INTO transactions (
		product_id, user_id, 
		type, quantity_change
	) VALUES ($1, $2, $3, $4)
	`

	_, err = tx.Exec(
		context.Background(),
		trxQuery,
		productId, 
		userId, 
		"IN", 
		req.Quantity, 
	)
	
	if err != nil {
		return Products{}, err
	}

	selectQuery := `
	SELECT 
		p.id, p.code_product, p.name, 
		p.image_url, p.purchase_price, 
		p.selling_price, p.quantity,
		u.name AS user,
		pc.name AS category
	FROM products p
	JOIN users u ON u.id = p.user_id
	JOIN product_categories pc ON pc.id = p.category_id
	WHERE p.id = $1
	`

	row, err := tx.Query(
		context.Background(),
		selectQuery,
		productId,
	)

	if err != nil {
		return Products{}, err
	}

	product, err := pgx.CollectOneRow[Products](row, pgx.RowToStructByName)
	if err != nil {
		return Products{}, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return Products{}, err
	}

	return product, nil
}