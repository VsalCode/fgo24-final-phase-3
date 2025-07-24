package controllers

import (
	"nashta_inventory/models"
	"nashta_inventory/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	req := models.RegisterRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Succcess: false,
			Message:  "Invalid Request!",
		})
		return
	}

	err = models.CreateNewUser(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Succcess: false,
			Message:  "Failed to register",
			Errors:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, utils.Response{
		Succcess: true,
		Message:  "Register Successfully!",
	})
}

func Login(ctx *gin.Context) {
	req := models.LoginRequest{}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{
			Succcess: false,
			Message:  "Invalid Request!",
		})
		return
	}

	result, err := models.ValidateLogin(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.Response{
			Succcess:false,
			Message: "Failed to login!",
			Errors: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Succcess: true,
		Message:  "Login Successfully!",
		Results:  result,
	})
}
