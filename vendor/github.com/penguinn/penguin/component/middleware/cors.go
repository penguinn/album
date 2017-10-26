package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/penguinn/penguin/component/config"
	"github.com/penguinn/penguin/utils"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var origin string
		origins := c.Request.Header["origin"]
		if len(origins) != 0 {
			origin = origins[0]
		}
		if config.Get("server.mode") == "debug" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			configOrigins := config.GetStringSlice("server.origin")
			ok := utils.InSlice(origin, configOrigins)
			if ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
	}
}
