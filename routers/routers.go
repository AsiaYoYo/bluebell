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

	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controllers.SignUpHandler)
	// 登录业务路由
	v1.POST("/login", controllers.LoginHandler)

	v1.GET("/version", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})
	return r
}
