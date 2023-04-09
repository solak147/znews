package controller

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
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

// @Summary 註冊 Step1
// @Tags user
// @version 1.0
// @produce application/json
// @param register body model.RegisterStep1 true "檢查帳號是否已存在"
// @Success 200 boolean successful return boolean
// @Router /member/registerStep1 [post]
func (u UsersController) CheckUserExit(c *gin.Context) {
	var form model.RegisterStep1

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	if service.CheckUserExit(form.Email) {

		// 設定隨機種子
		rand.Seed(time.Now().UnixNano())

		// 產生 6 個隨機數字
		var validCode string
		for len(validCode) < 6 {
			num := rand.Intn(9) + 1 // 產生 1 到 9 之間的隨機整數
			validCode += strconv.Itoa(num)
		}

		service.Send(validCode+"是您的驗證碼", "請在網頁上輸入您的驗證碼，此為系統自動發送的信件。", form.Email)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": validCode,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Account has exit",
			"data":   nil,
		})
	}

}

// @Summary 註冊 Step3
// @Tags user
// @version 1.0
// @produce application/json
// @param register body model.RegisterStep3 true "註冊帳號"
// @Success 200 boolean successful return boolean
// @Router /member/registerStep3 [post]
func (u UsersController) Register(c *gin.Context) {
	var form model.RegisterStep3
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	err := service.Register(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Register fail : " + err.Error(),
			"data":   nil,
		})
	} else {
		tokenString, _ := middleware.GenToken(form.Account)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{"token": tokenString},
		})
	}

}

// @Summary 登入
// @Tags user
// @version 1.0
// @produce application/json
// @param login body model.Login true "登入成功回傳 token"
// @Success 200 string successful return token
// @Router /member/login [post]
func (u UsersController) Login(c *gin.Context) {
	var form model.Login
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	user, err := service.GetUserByPwd(form.Account, form.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "Failed to parse params" + err.Error(),
			"data":   nil,
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": -1,
			"msg":    "User not found",
			"data":   nil,
		})

	} else {
		tokenString, _ := middleware.GenToken(form.Account)
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{"token": tokenString},
		})
	}

}

// t := gi18n.New()
// 	lan := c.Request.Header.Get("language")
// 	if lan == "" {
// 		lan = "zh"
// 	}
// 	t.SetLanguage(lan)
// 	t.Translate(c, "Response_Success")

// @Summary GetProfile
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param id path string true "帳號" default(”)
// @Success 200 string string successful return data
// @Router /member/profile/{id} [get]
func (u UsersController) GetProfile(c *gin.Context) {
	account := c.Params.ByName("account")

	user, err := service.GetUser(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    "User not found" + err.Error(),
			"data":   nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "Successfully get user data",
			"user":   &user,
		})
	}
}

// @Summary 個人資料儲存
// @Tags user
// @version 1.0
// @produce application/json
// @param MyAccount body model.ProfileSave true "修改成功回傳 boolean"
// @Success 200 boolean successful return boolean
// @Router /profile/save [post]
func (u UsersController) UpdateProfile(c *gin.Context) {
	var form model.ProfileSave
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	if err := service.UpdateUser(form); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status": -1,
			"msg":    err.Error(),
			"data":   false,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "Successfully save",
			"data":   true,
		})
	}

}
