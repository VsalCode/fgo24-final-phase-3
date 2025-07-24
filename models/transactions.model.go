package models

import (
	"context"
	"nashta_inventory/db"
)

func RecordTransaction(productId, userId int, transactionType string, quantityChange int) error {
	conn, err := db.DBConnect()
	if err != nil {
		return err
	}
	defer conn.Close()

	query := `
	INSERT INTO transactions (product_id, user_id, type, quantity_change)
	VALUES ($1, $2, $3, $4)
	`

	_, err = conn.Exec(
		context.Background(),
		query,
		productId, userId, transactionType, quantityChange)

	return err
}
