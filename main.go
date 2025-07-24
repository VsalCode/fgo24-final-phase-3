package main

import (
	"fmt"
	"nashta_inventory/db"
	"nashta_inventory/routers"
	"nashta_inventory/seeders"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.DBConnect()
	if err != nil {
		fmt.Println("Failed to connect database")
	}

	defer db.Close()

	err = seeders.SeedCategories(db)
	if err != nil {
		fmt.Println("Failed To seed product_categories")
	}

	r := gin.Default()
	routers.CombineRouters(r)

	r.Run(":8080")	
}