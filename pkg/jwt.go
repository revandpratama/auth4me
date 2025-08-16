package pkg

import (
	"errors"
	"strconv"
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

	expirationSecond, err := strconv.Atoi(config.ENV.JWT_EXPIRATION_SECOND)
	if err != nil || expirationSecond == 0 {
		expirationSecond = 30
	}

	expirationTime := time.Now().Add(time.Second * time.Duration(expirationSecond))

	claims := &CustomClaims{
		Email:        user.Email,
		Role:         user.Role,
		UserID:       user.ID,
		MFACompleted: mfaCompleted,
		Provider:     provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
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

func ParseExpiredToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.ENV.JWT_SECRET), nil
	}, jwt.WithoutClaimsValidation())

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}

	// Manually validate expiration if needed
	// For refresh flow, we expect it to be expired
	if claims.ExpiresAt == nil {
		return nil, errors.New("missing exp in token")
	}

	return claims, nil
}
