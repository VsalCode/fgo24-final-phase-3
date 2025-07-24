package controllers

import (
	"fmt"
	"nashta_inventory/dto"
	"nashta_inventory/models"
	"nashta_inventory/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GoodsMovement(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")

	if userId == "" && !exists {
		ctx.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized!",
		})
		return
	}

	req := dto.TransactionsRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Success: false,
			Message: "Invalid Request!",
			Errors: err.Error(),
		})
		return
	}

	result, err := models.AddNewTransactions(req, userId.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Transactions Failed!",
			Errors: err.Error(),
		})
		return
	}

	status := "incoming goods"
	if strings.ToLower(req.Type) == "out" {
		status = "outgoing goods"
	}

	ctx.JSON(http.StatusCreated, utils.Response{
		Success: true,
		Message: fmt.Sprintf("successfully recorded %s", status ),
		Results: result,
	})
}
