package dto

type ProductRequest struct {
	Name string `json:"name"`
	ImageURL string `json:"imageUrl" db:"image_url"`
	PurchasePrice float64 `json:"purchasePrice" db:"purchase_price"`
	SellingPrice float64 `json:"sellingPrice" db:"selling_price"`
	Quantity int `json:"quantity"`
	UserId int `json:"userId" db:"user_id"`
	CategoryId int `json:"categoryId" db:"category_id"`
}
