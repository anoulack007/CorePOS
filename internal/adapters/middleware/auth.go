package middleware

import (
	"net/http"
	"strings"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/pkg"
	"github.com/anoulack007/core-pos/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ContextUserIDKey  = "user_id"
	ContextStoreIDKey = "store_id"
	ContextRoleKey    = "role"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			pkg.Error(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			pkg.Error(c, http.StatusUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
		if tokenString == "" {
			pkg.Error(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			pkg.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			pkg.Error(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		// userID, err := util.GetUUIDClaim(claims, ContextUserIDKey)
		// if err != nil {
		// 	pkg.Error(c, http.StatusUnauthorized, "invalid user_id in token")
		// 	c.Abort()
		// 	return
		// }

		storeID, err := util.GetUUIDClaim(claims, ContextStoreIDKey)
		if err != nil {
			pkg.Error(c, http.StatusUnauthorized, "invalid store_id in token")
			c.Abort()
			return
		}

		roleValue, ok := claims[ContextRoleKey].(string)
		if !ok || strings.TrimSpace(roleValue) == "" {
			pkg.Error(c, http.StatusUnauthorized, "invalid role in token")
			c.Abort()
			return
		}

		role := domain.UserRole(roleValue)
		if !role.IsValid() {
			pkg.Error(c, http.StatusUnauthorized, "invalid role in token")
			c.Abort()
			return
		}

		c.Set(ContextRoleKey, role)
		c.Set(ContextStoreIDKey, storeID)
		c.Set(ContextRoleKey, role)
		c.Next()
	}
}

func AuthorizeStoreAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		storeID, exists := c.Get(ContextStoreIDKey)
		if !exists {
			pkg.Error(c, http.StatusForbidden, "store access denied")
			c.Abort()
			return
		}

		tokenStoreID, ok := storeID.(uuid.UUID)
		if !ok {
			pkg.Error(c, http.StatusForbidden, "store access denied")
			c.Abort()
			return
		}

		paramStoreID, err := uuid.Parse(c.Param("storeId"))
		if err != nil {
			pkg.Error(c, http.StatusBadRequest, "Invalid store ID")
			c.Abort()
			return
		}

		if tokenStoreID != paramStoreID {
			pkg.Error(c, http.StatusForbidden, "store access denied")
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireRoles(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	allowed := make(map[domain.UserRole]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		roleValue, exists := c.Get(ContextRoleKey)
		if !exists {
			pkg.Error(c, http.StatusForbidden, "role access denied")
			c.Abort()
			return
		}

		role, ok := roleValue.(domain.UserRole)
		if !ok {
			pkg.Error(c, http.StatusForbidden, "role access denied")
			c.Abort()
			return
		}

		if _, ok := allowed[role]; !ok {
			pkg.Error(c, http.StatusForbidden, "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
