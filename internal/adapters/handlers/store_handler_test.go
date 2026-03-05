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

func setupStoreRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/stores", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true, "data": body})
	})
	return r
}

func TestStoreCreate_200(t *testing.T) {
	r := setupStoreRouter()

	body := fmt.Sprintf(`{"name":"%s","plan_type":"%s"}`,
		gofakeit.Company(),
		gofakeit.RandomString([]string{"free", "pro", "enterprise"}),
	)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/stores", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestStoreCreate_400(t *testing.T) {
	r := setupStoreRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/stores", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
