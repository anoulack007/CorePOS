package services

import (
	"errors"
	"time"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo ports.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo ports.UserRepository, jwtSecret string) ports.AuthService{
	return &authService{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(user *domain.User, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return s.userRepo.Create(user)
}

func (s *authService) Login(username, password string) (string, error){
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"store_id": user.StoreID.String(),
		"username": user.Username,
		"role": user.Role,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}