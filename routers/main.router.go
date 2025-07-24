package routers

import (
	"github.com/gin-gonic/gin"
)

func CombineRouters(r *gin.Engine){
	authRouters(r.Group("/auth"))
	productsRouter(r.Group("/products"))
	transactionsRouter(r.Group("/transactions"))
}