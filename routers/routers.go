package routers

import (
	"net/http"
	"web_app/logger"
	"web_app/settings"

	"github.com/gin-gonic/gin"
)

// Setup 初始化routers
func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})
	return r
}
