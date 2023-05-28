package config

import (
	_ "znews/docs"

	"znews/app/controller"
	"znews/app/middleware"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CustomRouter(r *gin.Engine, m *persist.RedisStore) {

	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CorsMiddleware())

	//r.Use(middleware.ErrorMiddleware) 請求進來時會先初始化
	//r.Use(middleware.ErrorMiddleware()) 專案啟動時就初始化
	r.Use(middleware.ErrorMiddleware)

	r.GET("/ws", controller.SocketController().Socket)

	r.POST("/registerStep1", controller.UserController().CheckUserExit)
	r.POST("/registerStep3", controller.UserController().Register)
	r.POST("/login", controller.UserController().Login)

	r.GET("/case/get", controller.CaseController().GetCase)
	r.GET("/case/getDetail/:caseId", controller.CaseController().GetCaseDetail)

	// test
	// r.GET("/panic", func(c *gin.Context) {
	// 	panic("Something went wrong!")
	// })

	//r.GET("/image/:name", serveImage)

	r.Use(middleware.JWTAuthMiddleware())

	member := r.Group("/member")
	{
		member.GET("/profile/:account", controller.UserController().GetProfile)
		member.GET("/sohoSettingInit", controller.UserController().SohoSettingInit)
		member.POST("/profile/save", controller.UserController().UpdateProfile)
		member.POST("/sohoSetting", controller.UserController().SohoSetting)
		member.GET("/chkSohoSetting", controller.UserController().ChkSohoSetting)
		member.POST("/sohoUrl", controller.UserController().AddSohoUrl)
		member.GET("/sohoUrl", controller.UserController().GetSohoUrl)
		member.DELETE("/sohoUrl", controller.UserController().DeleteSohoUrl)

		//member.GET("/:id", middleware.JWTAuthMiddleware(), cache.CacheByRequestURI(m, 2*time.Hour), controller.UserController().GetUser)
	}

	casem := r.Group("/case")
	{
		casem.POST("/create", controller.CaseController().CreateCase)
		casem.POST("/quote", controller.CaseController().Quote)
		casem.GET("/quoteRecord/:deal", controller.CaseController().QuoteRecord)
		casem.GET("/getDetailAuth/:caseId", controller.CaseController().GetCaseDetail) // 找案件主檔及檔案
		casem.GET("/getDetailAuthOri/:caseId", controller.CaseController().GetCaseDetailOri)
		casem.POST("/update/:caseId", controller.CaseController().UpdateCase)
		casem.GET("/getFlow/:caseId", controller.CaseController().GetFlow)
		casem.POST("/flow", controller.CaseController().Flow)
		casem.POST("/collect", controller.CaseController().UpdateCollect)
		casem.GET("/collect", controller.CaseController().GetCollect)
		casem.GET("/release", controller.CaseController().GetCaseRelease)
	}

	file := r.Group("/file")
	{
		file.POST("/upload/:param", controller.FileController().Upload)
		//file.POST("/uploads", controller.FileController().Uploads)
		file.GET("/download/:caseId/:filename", controller.FileController().Download)
		file.GET("/sohoDownload/:filename/:param", controller.FileController().SohoDownload)
		file.GET("/sohowork/:param", controller.FileController().GetSohoWork)
		file.DELETE("/sohowork/:filename/:param", controller.FileController().DeleteSohoWork)
	}

	msg := r.Group("/message")
	{
		msg.GET("", controller.MsgController().GetMsgRecord)
		msg.GET("/chkNoRead", controller.MsgController().ChkNoRead)
		msg.GET("/:toAccount", controller.MsgController().GetMsgRecordDetail)
		msg.POST("/send", controller.MsgController().SendMsg)
		msg.POST("/deal", controller.MsgController().Deal)
		msg.PUT("/updateRead", controller.MsgController().UpdateRead)

	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

}

// 直接存取圖片
// func serveImage(c *gin.Context) {
// 	name := c.Params.ByName("name")
// 	imagePath := "app/service/images/" + name // 指定圖片檔案的路徑

// 	image, err := ioutil.ReadFile(imagePath)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Internal Server Error")
// 		return
// 	}

// 	c.Data(http.StatusOK, "image/jpeg", image)
// }
