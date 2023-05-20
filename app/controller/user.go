package controller

import (
	"fmt"
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

	if !service.CheckUserExist(form.Email) {

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
			"code": -1,
			"msg":  "Account has exit",
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
			"code": -1,
			"msg":  "Register fail : " + err.Error(),
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
			"code": -1,
			"msg":  "Failed to parse params" + err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": -1,
			"msg":  "User not found",
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

// @Summary SohoSetting
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param soho body string true "接案設定"
// @Success 200 string json successful return data
// @Router /member/sohoSetting [post]
func (u UsersController) SohoSetting(c *gin.Context) {
	account, _ := c.Get("account")

	var form model.SohoSettingForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Form bind error": err.Error()})
		return
	}

	if err := service.SohoSetting(fmt.Sprintf("%v", account), form); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Successfully save",
		})
	}
}

// @Summary 是否已填寫接案設定
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param soho path string true "是否已填寫接案設定"
// @Success 200 string json successful return data
// @Router /member/chkSohoSetting [get]
func (u UsersController) ChkSohoSetting(c *gin.Context) {
	account, _ := c.Get("account")

	if err := service.ChkSohoSetting(fmt.Sprintf("%v", account)); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Successfully save",
		})
	}
}

// @Summary SohoSetting
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param soho path string true "接案設定初始值"
// @Success 200 string json successful return data
// @Router /member/sohoSettingInit [get]
func (u UsersController) SohoSettingInit(c *gin.Context) {
	account, _ := c.Get("account")
	accountPath := c.Query("account")

	// x為自己看預覽檔案 其他為報價通知代入帳號查詢
	var tmp string
	if accountPath == "" || accountPath == "x" {
		tmp = fmt.Sprintf("%v", account)
	} else {
		tmp = accountPath
	}

	if data, err := service.SohoSettingInit(tmp); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
		})
	}
}

// @Summary 新增作品網址
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param soho body string true "新增作品網址"
// @Success 200 string json successful return data
// @Router /member/sohoUrl [post]
func (u UsersController) AddSohoUrl(c *gin.Context) {
	account, _ := c.Get("account")

	var form model.SohoUrlForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Form bind error" + err.Error(),
		})
		return
	}

	if err := service.AddSohoUrl(fmt.Sprintf("%v", account), form); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}
}

// @Summary 取得作品網址
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param soho path string true "取得作品網址"
// @Success 200 string json successful return data
// @Router /member/sohoUrl [get]
func (u UsersController) GetSohoUrl(c *gin.Context) {
	account, _ := c.Get("account")
	accountPath := c.Query("account")

	var tmp string
	if accountPath == "" || accountPath == "x" {
		tmp = fmt.Sprintf("%v", account)
	} else {
		tmp = accountPath
	}

	if data, err := service.GetSohoUrl(tmp); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": data,
		})
	}
}

// @Summary 刪除作品網址
// @Tags user
// @version 1.0
// @produce application/json
// @Security BearerAuth
// @param soho path string true "刪除作品網址"
// @Success 200 string json successful return data
// @Router /member/sohoUrl/{url} [delete]
func (u UsersController) DeleteSohoUrl(c *gin.Context) {
	account, _ := c.Get("account")

	var form model.SohoUrlForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "Form bind error" + err.Error(),
		})
		return
	}

	if err := service.DeleteSohoUrl(fmt.Sprintf("%v", account), form.Url); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "Success",
		})
	}
}
