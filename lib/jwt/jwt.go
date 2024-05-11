package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// JWTManager is a manager to handle JWT token
type JWTManager struct {
	secretKey     string
	tokenDuration int64
}

// NewJWTManager creates a new JWTManager
func NewJWTManager(secretKey string, tokenDuration int64) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Generate generates a new JWT token
func (manager *JWTManager) Generate(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// Verify verifies the JWT token
func (manager *JWTManager) Verify(accessToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(manager.secretKey), nil
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
func (manager *JWTManager) Refresh(accessToken string) (string, error) {
	claims, err := manager.Verify(accessToken)
	if err != nil {
		return "", err
	}
	return manager.Generate(claims)
}

// ExtractClaims extracts the claims from JWT token
func (manager *JWTManager) ExtractClaims(accessToken string) (map[string]interface{}, error) {
	claims, err := manager.Verify(accessToken)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ExtractClaimsFromRefreshToken extracts the claims from refresh token
func (manager *JWTManager) ExtractClaimsFromRefreshToken(refreshToken string) (map[string]interface{}, error) {
	claims, err := manager.Verify(refreshToken)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ExtractUserID extracts the user ID from JWT token
func (manager *JWTManager) ExtractUserID(accessToken string) (string, error) {
	claims, err := manager.ExtractClaims(accessToken)
	if err != nil {
		return "", err
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found")
	}
	return userID, nil
}

// ExtractUserIDFromRefreshToken extracts the user ID from refresh token
func (manager *JWTManager) ExtractUserIDFromRefreshToken(refreshToken string) (string, error) {
	claims, err := manager.ExtractClaimsFromRefreshToken(refreshToken)
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
func (manager *JWTManager) ExtractRole(accessToken string) (string, error) {
	claims, err := manager.ExtractClaims(accessToken)
	if err != nil {
		return "", err
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found")
	}
	return role, nil
}

// ExtractRoleFromRefreshToken extracts the role from refresh token
func (manager *JWTManager) ExtractRoleFromRefreshToken(refreshToken string) (string, error) {
	claims, err := manager.ExtractClaimsFromRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", fmt.Errorf("role not found")
	}
	return role, nil
}
