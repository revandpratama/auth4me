package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/revandpratama/auth4me/config"
	"github.com/revandpratama/auth4me/internal/auth/entity"
)

type CustomClaims struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	Provider     string `json:"provider,omitempty"` // Optional if OAuth
	SessionID    string `json:"sid,omitempty"`      // Optional, for token tracking
	MFACompleted bool   `json:"mfa,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entity.User, provider string, mfaCompleted bool) (string, error) {
	claims := &CustomClaims{
		Email: user.Email,
		Role:  user.Role,
		UserID: user.ID,
		MFACompleted: mfaCompleted,
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.ENV.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.ENV.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
