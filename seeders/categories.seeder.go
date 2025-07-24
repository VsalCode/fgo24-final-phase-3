package seeders

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool" 
)

func SeedCategories(pool *pgxpool.Pool) error { 
	ctx := context.Background()

	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM product_categories").Scan(&count) 
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Println("Product categories already exist. Skipping seeding.")
		return nil
	}

	categories := []struct {
		Name        string
		Description string
	}{
		{Name: "elektronik", Description: "Berbagai gadget dan perangkat elektronik."},
		{Name: "makanan", Description: "Berbagai jenis produk makanan siap saji dan bahan mentah."},
		{Name: "minuman", Description: "Aneka minuman seperti air mineral, jus, kopi, teh, dan minuman bersoda."},
		{Name: "pakaian", Description: "Busana dan aksesoris untuk pria, wanita, dan anak-anak."},
		{Name: "perabotan", Description: "Peralatan dapur, dan perabot rumah tangga lainnya."},
	}

	for _, category := range categories {
		_, err := pool.Exec(ctx, 
			`INSERT INTO product_categories (name, description) VALUES ($1, $2)`,
			category.Name, category.Description)
		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully seeded product categories.")
	return nil
}
