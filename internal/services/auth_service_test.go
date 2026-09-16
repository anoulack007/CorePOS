package services

import (
	"testing"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUserRepositoryStub struct {
	user         *domain.User
	createdOwner *domain.User
}

func (r *authUserRepositoryStub) FindByID(storeID, id uuid.UUID) (*domain.User, error) {
	return r.user, nil
}

func (r *authUserRepositoryStub) FindByUsername(username string) (*domain.User, error) {
	return r.user, nil
}

func (r *authUserRepositoryStub) Create(user *domain.User) error {
	return nil
}

func (r *authUserRepositoryStub) CreateInitialOwner(user *domain.User) error {
	r.createdOwner = user
	return nil
}

func TestRegisterForcesInitialUserToOwner(t *testing.T) {
	repo := &authUserRepositoryStub{}
	service := NewAuthService(repo, "test-secret")
	user := &domain.User{Role: domain.RoleAdmin}

	if err := service.Register(user, "strong-password"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if repo.createdOwner == nil || repo.createdOwner.Role != domain.RoleOwner {
		t.Fatalf("expected initial user role %q", domain.RoleOwner)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("strong-password")) != nil {
		t.Fatal("expected password to be stored as a bcrypt hash")
	}
}

func TestLoginCreatesDistinctAccessAndRefreshTokenTypes(t *testing.T) {
	password := "strong-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	repo := &authUserRepositoryStub{user: &domain.User{
		ID:           uuid.New(),
		StoreID:      uuid.New(),
		Username:     "owner",
		PasswordHash: string(hash),
		Role:         domain.RoleOwner,
	}}
	service := NewAuthService(repo, "test-secret")

	accessToken, refreshToken, err := service.Login("owner", password)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if got := parsedTokenType(t, accessToken, "test-secret"); got != "access" {
		t.Fatalf("access token type = %q", got)
	}
	if got := parsedTokenType(t, refreshToken, "test-secret"); got != "refresh" {
		t.Fatalf("refresh token type = %q", got)
	}

	if _, _, err := service.RefreshToken(accessToken); err == nil {
		t.Fatal("expected access token to be rejected by RefreshToken")
	}
}

func parsedTokenType(t *testing.T, tokenString, secret string) string {
	t.Helper()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("failed to parse generated token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("generated token did not contain map claims")
	}
	tokenType, _ := claims["token_type"].(string)
	return tokenType
}
