package middleware

import (
	"log"
	"net/http"

	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from panics.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("\033[31m[PANIC RECOVERED]\033[0m %v", err)
				pkg.Error(c, http.StatusInternalServerError, "Internal Server Error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
