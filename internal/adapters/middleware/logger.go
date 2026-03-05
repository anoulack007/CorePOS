package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a clean, readable request logger middleware.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method

		// Read body
		var bodyStr string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			bodyStr = string(bodyBytes)
			// Restore body for handlers
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process request
		c.Next()

		// After request
		duration := time.Since(start)
		status := c.Writer.Status()
		reqID, _ := c.Get(RequestIDKey)

		statusColor := colorForStatus(status)
		methodColor := colorForMethod(method)
		reset := "\033[0m"

		// Build full path
		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		// Truncate body if too long
		if len(bodyStr) > 200 {
			bodyStr = bodyStr[:200] + "..."
		}

		// Log format
		log.Printf("%s %3d %s | %10v | %-15s | %s %-7s %s %s",
			statusColor, status, reset,
			duration.Round(time.Microsecond),
			c.ClientIP(),
			methodColor, method, reset,
			fullPath,
		)

		// Log request ID + body (if present)
		if reqID != nil {
			fmt.Printf("         └─ reqID: %s\n", reqID)
		}
		if bodyStr != "" && bodyStr != "{}" {
			fmt.Printf("         └─ body:  %s\n", bodyStr)
		}
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[42;30m" // Green bg, black text
	case code >= 300 && code < 400:
		return "\033[43;30m" // Yellow bg, black text
	case code >= 400 && code < 500:
		return "\033[41;97m" // Red bg, white text
	default:
		return "\033[45;97m" // Magenta bg, white text
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
