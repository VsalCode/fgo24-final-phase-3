package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DBConnect() (*pgxpool.Pool, error ){
	connectionString := "postgres://faisal:faisalg5@localhost:5432/nashta_inventory"

	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	return pool, nil
}