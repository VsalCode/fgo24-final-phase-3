package models

import (
	"context"
	"time"
	"nashta_inventory/db"

	"github.com/jackc/pgx/v5"
)

type TransactionsRequest struct {
	ProductID      int    `json:"productId" db:"product_id" binding:"required"`
	Type           string `json:"type" binding:"required"`
	QuantityChange int    `json:"quantityChange" db:"quantity_change" binding:"required"`
}

type TransactionsResponse struct {
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

func AddNewTransactions(req TransactionsRequest, userId int) (*TransactionsResponse, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var trxId int
	insertQuery := `
		INSERT INTO transactions (product_id, user_id, type, quantity_change)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = tx.QueryRow(
		context.Background(),
		insertQuery,
		req.ProductID, userId, req.Type, req.QuantityChange).Scan(&trxId)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	_, err = tx.Exec(context.Background(), updateQuery, req.QuantityChange, req.ProductID)
	if err != nil {
		return nil, err
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

	rows, err := tx.Query(context.Background(), joinQuery, trxId)
	if err != nil {
		return nil, err
	}

	result, err := pgx.CollectOneRow[TransactionsResponse](rows, pgx.RowToStructByName)

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetTransactionHistory(productId int, limit, offset int) ([]*TransactionsResponse, error) {
	conn, err := db.DBConnect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `
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
		WHERE t.product_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := conn.Query(context.Background(), query, productId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*TransactionsResponse
	for rows.Next() {
		var t TransactionsResponse
		err := rows.Scan(
			&t.ID,
			&t.ProductName,
			&t.CategoryName,
			&t.Type,
			&t.QuantityChange,
			&t.PurchasePrice,
			&t.SellingPrice,
			&t.Stock,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	return transactions, nil
}