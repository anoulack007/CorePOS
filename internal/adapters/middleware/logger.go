package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Logger returns a clean, readable request logger middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		// Process request
		c.Next()
		// After request
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		// Color based on status
		statusColor := colorForStatus(status)
		methodColor := colorForMethod(method)
		resetColor := "\033[0m"
		log.Printf("| %s %3d %s | %13v | %-15s | %s %-7s %s %s",
			statusColor, status, resetColor,
			duration,
			c.ClientIP(),
			methodColor, method, resetColor,
			path,
		)
	}
}


func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[42m" // Green
	case code >= 300 && code < 400:
		return "\033[43m" // Yellow
	case code >= 400 && code < 500:
		return "\033[41m" // Red
	default:
		return "\033[45m" // Magenta
	}
}
func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // Blue
	case "POST":
		return "\033[32m" // Green
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[36m" // Cyan
	default:
		return "\033[0m"
	}
}