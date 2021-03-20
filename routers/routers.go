package routers

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/settings"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup 初始化routers
func Setup() *gin.Engine {
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.POST("/signup", controllers.SignUpHandler)

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})
	return r
}
