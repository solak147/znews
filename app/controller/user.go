package controller

import (
	"net/http"
	"strconv"
	"znews/app/middleware"
	"znews/app/model"
	"znews/app/service"

	"github.com/gin-gonic/gin"
)

type UsersController struct{}

func UserController() UsersController {
	return UsersController{}
}

type Register struct {
	Account  string `json:"account" binding:"required" example:"account"`
	Password string `json:"password" binding:"required" example:"password"`
	Email    string `json:"email" binding:"required" example:"test123@gmail.com"`
}

func (u UsersController) AuthHandler(c *gin.Context) {
	var form model.User
	bindErr := c.BindJSON(&form)
	if bindErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Invalid params",
		})
		return
	}

	userOne, err := service.LoginOneUser(form.Account, form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params" + err.Error(),
			"data":   nil,
		})
		return
	}

	if userOne == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg":    "User not found",
			"data":   nil,
		})
		return
	}

	if userOne != nil {
		tokenString, _ := middleware.GenToken(form.Account)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{"token": tokenString},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Verified Failed.",
	})

}

func (u UsersController) CreateUser(c *gin.Context) {
	var form Register
	bindErr := c.BindJSON(&form)

	if bindErr == nil {
		err := service.RegisterOneUser(form.Account, form.Password, form.Email)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 1,
				"msg":    "success Register",
				"data":   nil,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": -1,
				"msg":    "Register Failed" + err.Error(),
				"data":   nil,
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse register data" + bindErr.Error(),
			"data":   nil,
		})
	}
}

func (u UsersController) GetUser(c *gin.Context) {
	id := c.Params.ByName("id")

	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params" + err.Error(),
			"data":   nil,
		})
	}

	userOne, err := service.SelectOneUsers(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg":    "User not found" + err.Error(),
			"data":   nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "Successfully get user data",
			"user":   &userOne,
		})
	}
}
