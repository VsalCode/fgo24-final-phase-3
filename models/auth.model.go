package models

import (
	"context"
	"nashta_inventory/db"
	"nashta_inventory/utils"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" db:"password_hash"`
	Phone    string `json:"phone"`
}

func CreateNewUser(req RegisterRequest) error {
	conn, err := db.DBConnect()
	if err != nil {
		return err
	}

	defer conn.Close()

	query := `
	INSERT INTO users (name, email, password_hash, phone)
	VALUES ($1, $2, $3, $4)
	`

	password, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	_, err = conn.Exec(
		context.Background(),
		query,
		req.Name,
		req.Email,
		password ,
		req.Phone,
	)

	if err != nil {
		return err
	}

	return nil
}	
