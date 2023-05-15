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

	r.Use(middleware.LoggerToFile())
	r.Use(corsMiddleware())

	r.GET("/ws", controller.SocketController().Socket)

	r.POST("/registerStep1", controller.UserController().CheckUserExit)
	r.POST("/registerStep3", controller.UserController().Register)
	r.POST("/login", controller.UserController().Login)

	r.GET("/case/get", controller.CaseController().GetCase)
	r.GET("/case/getDetail/:caseid", controller.CaseController().GetCaseDetail)

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
		casem.GET("/quoteRecord", controller.CaseController().QuoteRecord)
		casem.GET("/chkBefQuote/:caseId", controller.CaseController().ChkBefQuote)
		casem.GET("/getDetailAuth/:caseId", controller.CaseController().GetCaseDetail)

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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 1天

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// router.GET("/user/:name", func(c *gin.Context) {
// 	fmt.Println(c.FullPath())  // /user/:name/
// 	name := c.Param("name")
// 	c.String(http.StatusOK, "Hello %s", name)
// })

// server.GET("/hc", func(c *gin.Context) {
// 	c.JSON(200, gin.H{
// 		"message": "health check",
// 	})
// })
