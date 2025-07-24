package models

import (
	"context"
	"fmt"
	"nashta_inventory/db"
	"nashta_inventory/utils"
)

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" db:"password_hash" binding:"required,min=6"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	Phone           string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" db:"password_hash" binding:"required"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
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
		password,
		req.Phone,
	)

	if err != nil {
		return err
	}

	return nil
}

func ValidateLogin(req LoginRequest) (User, error) {
    conn, err := db.DBConnect()
    if err != nil {
        return User{}, err
    }
    defer conn.Close()

    var userPassword string
    err = conn.QueryRow(
        context.Background(),
        "SELECT password_hash FROM users WHERE email = $1",
        req.Email,
    ).Scan(&userPassword)

    if err != nil {
        return User{}, err
    }

    isValid := utils.CheckPasswordHash(req.Password, userPassword)
    if !isValid {
        return User{}, fmt.Errorf("invalid password")
    }

    var user User
    query := `
    SELECT id, name, email, phone FROM users
    WHERE email = $1
    `
    err = conn.QueryRow(context.Background(), query, req.Email).Scan(
        &user.Id, 
        &user.Name, 
        &user.Email, 
        &user.Phone,
    )

    if err != nil {
        return User{}, err
    }

    return user, nil
}
