package routers

import (
	"nashta_inventory/controllers"
	"nashta_inventory/middlewares"

	"github.com/gin-gonic/gin"
)

func productsRouter(r *gin.RouterGroup){
	r.Use(middlewares.VerifyToken())
	r.GET("/categories", controllers.GetAllCategories)
	r.GET("/", controllers.GetAllProducts )
}