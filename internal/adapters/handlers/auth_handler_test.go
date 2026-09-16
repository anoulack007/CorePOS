package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/gin-gonic/gin"
)

type authServiceStub struct {
	loginErr error
}

func (s *authServiceStub) Register(user *domain.User, password string) error { return nil }
func (s *authServiceStub) Login(username, password string) (string, string, error) {
	if s.loginErr != nil {
		return "", "", s.loginErr
	}
	return "access-token", "refresh-token", nil
}
func (s *authServiceStub) RefreshToken(token string) (string, string, error) {
	return "access-token", "refresh-token", nil
}
func (s *authServiceStub) Logout() error { return nil }

func TestAuthLoginUsesProductionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&authServiceStub{}, nil, nil)
	router := gin.New()
	router.POST("/auth/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"owner","password":"strong-password"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "access-token") {
		t.Fatalf("expected successful production login response, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAuthLoginRejectsInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&authServiceStub{loginErr: errors.New("invalid credentials")}, nil, nil)
	router := gin.New()
	router.POST("/auth/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"owner","password":"wrong-password"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}
