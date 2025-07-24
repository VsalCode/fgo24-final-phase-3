package models

import (
	"context"
	"fmt"
	"nashta_inventory/db"
	"nashta_inventory/dto"
	"time"

	"github.com/jackc/pgx/v5"
)

type Transactions struct {
	ID            int       `json:"id" db:"id"`
	ProductName   string    `json:"productName" db:"product_name"`
	CategoryName  string    `json:"categoryName" db:"category_name"`
	Type          string    `json:"type" db:"type"`
	QuantityChange int      `json:"quantityChange" db:"quantity_change"`
	PurchasePrice float64   `json:"purchasePrice" db:"purchase_price"`
	SellingPrice  float64   `json:"sellingPrice" db:"selling_price"`
	Stock         int       `json:"stock" db:"stock"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

func AddNewTransactions(req dto.TransactionsRequest, userId int) (*Transactions, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if req.Type == "out" || req.Type == "OUT" {
		var currentStock int
		stockQuery := `SELECT quantity FROM products WHERE id = $1`
		err = conn.QueryRow(context.Background(), stockQuery, req.ProductID).Scan(&currentStock)
		if err != nil {
			return nil, fmt.Errorf("product not found!")
		}

		if currentStock < req.QuantityChange {
			return nil, fmt.Errorf("stock: available %d, requested %d", currentStock, req.QuantityChange)
		}
	}

	var trxId int
	insertQuery := `
		INSERT INTO transactions (product_id, user_id, type, quantity_change)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = conn.QueryRow(
		context.Background(),
		insertQuery,
		req.ProductID, userId, req.Type, req.QuantityChange).Scan(&trxId)
	if err != nil {
		return nil, fmt.Errorf("failed to insert transaction: %v", err)
	}

	var updateQuery string
	switch req.Type {
	case "in", "IN":
		updateQuery = `
			UPDATE products 
			SET quantity = quantity + $1, updated_at = NOW()
			WHERE id = $2
		`
	case "out", "OUT":
		updateQuery = `
			UPDATE products 
			SET quantity = quantity - $1, updated_at = NOW()
			WHERE id = $2
		`
	default:
		return nil, fmt.Errorf("invalid transaction type: %s", req.Type)
	}

	_, err = conn.Exec(context.Background(), updateQuery, req.QuantityChange, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to update product quantity: %v", err)
	}

	joinQuery := `
		SELECT 
			t.id,
			p.name AS product_name,
			pc.name AS category_name,
			t.type,
			t.quantity_change,
			p.purchase_price,
			p.selling_price,
			p.quantity AS stock,
			t.created_at
		FROM transactions t
		JOIN products p ON p.id = t.product_id
		JOIN product_categories pc ON pc.id = p.category_id
		WHERE t.id = $1
	`

	rows, err := conn.Query(context.Background(), joinQuery, trxId)
	if err != nil {
		return nil, fmt.Errorf("failed to query transaction result: %v", err)
	}
	defer rows.Close()

	result, err := pgx.CollectOneRow[Transactions](rows, pgx.RowToStructByName)
	if err != nil {
		return nil, fmt.Errorf("failed to collect transaction result: %v", err)
	}

	return &result, nil
}
