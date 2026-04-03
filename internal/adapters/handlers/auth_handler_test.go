package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
)

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/auth/login", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true, "access_token": gofakeit.LetterN(32)})
	})
	return r
}

func TestAuthLogin_200(t *testing.T) {
	r := setupAuthRouter()
	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, gofakeit.Username(), gofakeit.Password(true, true, true, true, false, 8))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestAuthLogin_400(t *testing.T) {
	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
