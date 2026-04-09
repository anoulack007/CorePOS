package util

import (
	"github.com/anoulack007/core-pos/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NormalizeUserRole(role string) domain.UserRole {
	r := domain.UserRole(role)
	if !r.IsValid() {
		return domain.RoleCashier
	}
	return r
}

func GetUUIDClaim(claims jwt.MapClaims, key string) (uuid.UUID, error) {
	value, ok := claims[key].(string)
	if !ok {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	return uuid.Parse(value)
}
