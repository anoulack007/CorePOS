package middleware

import (
	"net/http"
	"strings"

	"github.com/anoulack007/core-pos/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(jwtSecret string) gin.HandlerFunc{
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer"){
			pkg.Error(c, http.StatusUnauthorized,"unauthorized")
			c.Abort();
			return 
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			pkg.Error(c, http.StatusUnauthorized,"invalid token")
			c.Abort();
			return 
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", claims["user_id"])
		c.Set("store_id", claims["store_id"])
		c.Next()
	}
}