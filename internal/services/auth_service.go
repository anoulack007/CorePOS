package services

import (
	"errors"
	"time"

	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/anoulack007/core-pos/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo  ports.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo ports.UserRepository, jwtSecret string) ports.AuthService {
	return &authService{
		userRepo:  userRepo,
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

func (s *authService) Login(username, password string) (string, string, error) { 
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid credentials!") 
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials") 
	}
	accessToken, _ := s.generateToken(user, 15*time.Minute)
	refreshToken, _ := s.generateToken(user, 7*24*time.Hour)
	return accessToken, refreshToken, nil 
}

func (s *authService) RefreshToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })
    if err != nil || !token.Valid {
        return "", "", errors.New("invalid refresh token")
    }
    claims := token.Claims.(jwt.MapClaims)

	storeID, err := uuid.Parse(claims["store_id"].(string))
	if err != nil {
		return "", "", errors.New("invalid store_id in token")
	}

	userID, err := uuid.Parse(claims["user_id"].(string))
	if err != nil {
		return "", "", errors.New("invalid user_id in token")
	}

	user, err := s.userRepo.FindByID(storeID, userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}


    accessToken, _ := s.generateToken(user, 15*time.Minute)
    refreshToken, _ := s.generateToken(user, 7*24*time.Hour)
    return accessToken, refreshToken, nil
}

func (s *authService) Logout() error {
    return nil
}


func (s *authService) generateToken(user *domain.User, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"store_id": user.StoreID.String(),
		"username": user.Username,
		"role": user.Role,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return  token.SignedString([]byte(s.jwtSecret))
}