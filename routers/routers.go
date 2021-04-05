package routers

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"bluebell/settings"
	"fmt"
	"net/http"

	_ "bluebell/docs"

	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

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

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		// 获取community列表
		v1.GET("/community", controllers.CommunityHandler)
		// 获取community列表
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		// 创建帖子路由
		v1.POST("/post", controllers.CreatePostHandler)
		// 获取单个帖子详情
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		// 获取帖子列表
		v1.GET("/posts", controllers.GetPostListHandler)
		// 根据path中的query参数获取帖子列表
		v1.GET("/posts2", controllers.GetPostListHandler2)

		// 为帖子投票
		v1.POST("/vote", controllers.PostVoteHandler)
	}

	return r
}
