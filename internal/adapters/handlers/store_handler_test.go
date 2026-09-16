package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStoreCreateRejectsBlankName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewStoreHandler(nil)
	router := gin.New()
	router.POST("/stores", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/stores", strings.NewReader(`{"name":"  "}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", w.Code, w.Body.String())
	}
}
