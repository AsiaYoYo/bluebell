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

	v1.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		// 获取community列表
		v1.GET("/community", controllers.CommunityHandler)
		// 获取community列表
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		// 创建帖子路由
		v1.POST("/post", controllers.CreatePostHandler)

		v1.GET("/version2", func(c *gin.Context) {
			c.String(http.StatusOK, settings.Conf.Version)
		})
	}

	return r
}
