package models

type TransactionsRequest struct {
	ProductId      int    `json:"productId" db:"product_id"`
	Type           string `json:"type"`
	QuantityChange int    `json:"quantityChange" db:"quantity_change"`
}
