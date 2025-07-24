package dto

type TransactionsRequest struct {
	ProductID      int    `json:"productId" db:"product_id" binding:"required"`
	Type           string `json:"type" binding:"required"`
	QuantityChange int    `json:"quantityChange" db:"quantity_change" binding:"required"`
}
