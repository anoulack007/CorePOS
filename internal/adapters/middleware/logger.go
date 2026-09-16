package middleware

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger records request metadata without logging request bodies or credentials.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		reqID, _ := c.Get(RequestIDKey)

		statusColor := colorForStatus(status)
		methodColor := colorForMethod(method)
		reset := "\033[0m"
		dim := "\033[90m"

		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		fmt.Printf("\n%s%s%s\n", dim, strings.Repeat("─", 80), reset)
		log.Printf("%s %3d %s │ %10v │ %-15s │ %s %-7s %s %s",
			statusColor, status, reset,
			duration.Round(time.Microsecond),
			c.ClientIP(),
			methodColor, method, reset,
			fullPath,
		)

		if reqID != nil {
			fmt.Printf("  %s└─ id: %s%s\n", dim, reset, reqID)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[42;30m"
	case code >= 300 && code < 400:
		return "\033[43;30m"
	case code >= 400 && code < 500:
		return "\033[41;97m"
	default:
		return "\033[45;97m"
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[36m"
	case "POST":
		return "\033[32m"
	case "PUT":
		return "\033[33m"
	case "DELETE":
		return "\033[31m"
	case "PATCH":
		return "\033[35m"
	default:
		return "\033[0m"
	}
}
