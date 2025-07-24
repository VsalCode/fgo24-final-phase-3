package routers

import (
	"nashta_inventory/controllers"
	"nashta_inventory/middlewares"

	"github.com/gin-gonic/gin"
)

func transactionsRouter(r *gin.RouterGroup){
	r.Use(middlewares.VerifyToken())
	r.POST("", controllers.GoodsMovement)
}