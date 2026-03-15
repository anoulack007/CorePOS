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

func setupCategoryRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/categories", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
	})
	return r
}

func TestCategoryCreate_200(t *testing.T) {
	r := setupCategoryRouter()
	body := fmt.Sprintf(`{"name":"%s"}`, gofakeit.Word())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCategoryCreate_400(t *testing.T) {
	r := setupCategoryRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
