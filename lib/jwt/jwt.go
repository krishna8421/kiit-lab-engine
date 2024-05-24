package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager is a manager to handle JWT token
type JWTManager struct {
	SecretKey            string
	AccessTokenDuration  int64
	RefreshTokenDuration int64
}

// NewJWTManager creates a new JWTManager
func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration int64) *JWTManager {
	return &JWTManager{secretKey, accessTokenDuration, refreshTokenDuration}
}

// Generate generates a new JWT token
func (manager *JWTManager) Generate(payload map[string]interface{}, isRefreshToken bool) (string, error) {
	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}
	if isRefreshToken {
		claims["exp"] = time.Now().Add(time.Duration(manager.RefreshTokenDuration) * time.Second).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Duration(manager.AccessTokenDuration) * time.Second).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.SecretKey))
}

// Verify verifies the JWT token
func (manager *JWTManager) Verify(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return claims, nil
}

// Refresh refreshes the JWT token
func (manager *JWTManager) Refresh(refreshToken string) (string, error) {
	claims, err := manager.Verify(refreshToken)
	if err != nil {
		return "", err
	}
	return manager.Generate(claims, false)
}

// ExtractClaims extracts the claims from JWT token
func (manager *JWTManager) ExtractClaims(tokenString string) (map[string]interface{}, error) {
	claims, err := manager.Verify(tokenString)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ExtractUserID extracts the user ID from JWT token
func (manager *JWTManager) ExtractUserID(tokenString string) (string, error) {
	claims, err := manager.ExtractClaims(tokenString)
	if err != nil {
		return "", err
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found")
	}
	return userID, nil
}

// ExtractRole extracts the role from JWT token
func (manager *JWTManager) ExtractRole(tokenString string) (string, error) {
	claims, err := manager.ExtractClaims(tokenString)
	if err != nil {
		return "", err
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found")
	}
	return role, nil
}
