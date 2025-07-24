package routers

import (
	"nashta_inventory/controllers"
	"github.com/gin-gonic/gin"
)

func authRouters(r *gin.RouterGroup){
	r.POST("/register", controllers.Register )
	r.POST("/login", controllers.Login)
}