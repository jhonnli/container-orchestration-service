package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonnli/logs"
)

const (
	APPID  = "appId"
	VER    = "ver"
	TIME   = "time"
	SIGN   = "sign"
	SECRET = "secret"

	VERSION = "1.0"
)

func ClientAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := c.GetHeader(SECRET)
		appId := c.GetHeader(APPID)
		if secret == "" || appId == "" {
			logs.Warn("获取%s时secret/appId校验失败", c.Request.URL)
			c.Abort()
			return
		}
		c.Next()
	}
}
