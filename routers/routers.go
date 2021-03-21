package routers

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup 初始化routers
func Setup(mode string) *gin.Engine {
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
	}
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 注册业务路由
	r.POST("/signup", controllers.SignUpHandler)
	// 登录业务路由
	r.POST("/login", controllers.LoginHandler)

	r.GET("/version", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})
	return r
}
