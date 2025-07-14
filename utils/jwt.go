package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      int          `json:"user_id"`
	Email       string       `json:"email"`
	Tipo        string       `json:"tipo"`
	RoleID      *int         `json:"role_id,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func GenerateAccessToken(userID int, email, tipo string, roleID *int, permissions []Permission) (string, error) {
	fmt.Printf("[JWT] 🔑 Generando token para usuario ID: %d, Email: %s\n", userID, email)
	if roleID != nil {
		fmt.Printf("[JWT] 👤 Role ID: %d\n", *roleID)
	}
	fmt.Printf("[JWT] 🔐 Permisos incluidos en token: %d\n", len(permissions))
	for _, perm := range permissions {
		fmt.Printf("[JWT]   - %s:%s\n", perm.Resource, perm.Action)
	}

	claims := Claims{
		UserID:      userID,
		Email:       email,
		Tipo:        tipo,
		RoleID:      roleID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), 
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "hospital-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		fmt.Printf("[JWT] ❌ JWT_SECRET no configurado\n")
		return "", errors.New("JWT_SECRET not set")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Printf("[JWT] ❌ Error al firmar token: %v\n", err)
		return "", err
	}

	fmt.Printf("[JWT] ✅ Token generado exitosamente\n")
	return tokenString, nil
}

func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   string(rune(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), 
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "hospital-system",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	return token.SignedString([]byte(jwtSecret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
