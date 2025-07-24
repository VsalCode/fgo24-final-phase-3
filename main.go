package main

import (
	"fmt"
	"nashta_inventory/db"
	"nashta_inventory/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.DBConnect()
	if err != nil {
		fmt.Println("Failed to connect database")
	}

	defer db.Close()

	r := gin.Default()
	routers.CombineRouters(r)

	r.Run(":8080")	
}