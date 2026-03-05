package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
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
			bodyStr = compactJSON(string(bodyBytes))
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
		dim := "\033[90m"

		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		// Truncate body
		if len(bodyStr) > 150 {
			bodyStr = bodyStr[:150] + "..."
		}

		// Main log line
		fmt.Printf("\n%s%s%s\n", dim, strings.Repeat("─", 80), reset)
		log.Printf("%s %3d %s │ %10v │ %-15s │ %s %-7s %s %s",
			statusColor, status, reset,
			duration.Round(time.Microsecond),
			c.ClientIP(),
			methodColor, method, reset,
			fullPath,
		)

		if reqID != nil {
			fmt.Printf("  %s├─ id:   %s%s\n", dim, reset, reqID)
		}
		if bodyStr != "" && bodyStr != "{}" {
			fmt.Printf("  %s└─ body: %s%s\n", dim, reset, bodyStr)
		}
	}
}

// compactJSON removes newlines and extra spaces from JSON string.
func compactJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\t", "")
	// Collapse multiple spaces
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[42;30m" // Green bg
	case code >= 300 && code < 400:
		return "\033[43;30m" // Yellow bg
	case code >= 400 && code < 500:
		return "\033[41;97m" // Red bg
	default:
		return "\033[45;97m" // Magenta bg
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[36m" // Cyan
	case "POST":
		return "\033[32m" // Green
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[35m" // Magenta
	default:
		return "\033[0m"
	}
}
