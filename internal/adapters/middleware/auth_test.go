package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestAuth_AllowsValidBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	storeID := uuid.New()
	userID := uuid.New()
	secret := "test-secret"

	r := gin.New()
	r.GET("/protected", Auth(secret), func(c *gin.Context) {
		if _, ok := c.Get(ContextUserIDKey); !ok {
			t.Fatalf("expected %s in context", ContextUserIDKey)
		}
		if _, ok := c.Get(ContextStoreIDKey); !ok {
			t.Fatalf("expected %s in context", ContextStoreIDKey)
		}
		if _, ok := c.Get(ContextRoleKey); !ok {
			t.Fatalf("expected %s in context", ContextRoleKey)
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken(t, secret, userID, storeID, "admin"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestAuthorizeStoreAccess_BlocksDifferentStore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	tokenStoreID := uuid.New()
	requestStoreID := uuid.New()
	userID := uuid.New()

	r := gin.New()
	r.GET("/stores/:storeId/products", Auth(secret), AuthorizeStoreAccess(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/stores/"+requestStoreID.String()+"/products", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken(t, secret, userID, tokenStoreID, "admin"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestRequireRoles_BlocksDisallowedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	storeID := uuid.New()
	userID := uuid.New()

	r := gin.New()
	r.GET("/admin", Auth(secret), RequireRoles("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken(t, secret, userID, storeID, "cashier"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestRequireRoles_AllowsAllowedRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test-secret"
	storeID := uuid.New()
	userID := uuid.New()

	r := gin.New()
	r.POST("/admin", Auth(secret), RequireRoles("owner", "admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req := httptest.NewRequest(http.MethodPost, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+signedToken(t, secret, userID, storeID, "admin"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func signedToken(t *testing.T, secret string, userID, storeID uuid.UUID, role string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		ContextUserIDKey:  userID.String(),
		ContextStoreIDKey: storeID.String(),
		ContextRoleKey:    role,
		"exp":             time.Now().Add(time.Hour).Unix(),
		"iat":             time.Now().Unix(),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signed
}
