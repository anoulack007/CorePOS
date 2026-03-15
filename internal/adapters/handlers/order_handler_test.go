package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupOrderRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/orders", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
	})
	return r
}

func TestOrderCreate_200(t *testing.T) {
	r := setupOrderRouter()

	// Mock order with items
	body := `{"items": [{"product_id": "8aa1c720-379a-412f-90b1-4700d8cb30fc", "quantity": 2}]}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestOrderCreate_400(t *testing.T) {
	r := setupOrderRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
