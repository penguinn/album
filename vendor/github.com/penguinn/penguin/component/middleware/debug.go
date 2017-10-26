package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/penguinn/penguin/component/log"
	"time"
)

func DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		log.Debugf(
			"start: %v	end: %v	 latency: %v  statusCode: %v  clientIP: %v  method: %v  path: %v",
			start, end, latency, statusCode, clientIP, method, path,
		)
	}
}
