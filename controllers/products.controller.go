package controllers

import (
	"nashta_inventory/models"
	"nashta_inventory/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCategories(ctx *gin.Context){
	userId, exists := ctx.Get("userId")

	if userId != "" && !exists {
		ctx.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized!",
		})
		return
	}

	result, err := models.FindAllCategories() 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed To Get All Categories",
			Errors: err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Get All Categories Successfully!",
		Results: result,
	})

}

func GetAllProducts(ctx *gin.Context){
	userId, exists := ctx.Get("userId")

	if userId != "" && !exists {
		ctx.JSON(http.StatusUnauthorized, utils.Response{
			Success: false,
			Message: "Unauthorized!",
		})
		return
	}

	result, err := models.FindAllProducts() 
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Success: false,
			Message: "Failed To Get All Products!",
			Errors: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Success: true,
		Message: "Get All Products Successfully!",
		Results: result,
	})

}