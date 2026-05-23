package utils

import (
	timePackage "time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pnaskardev/URL-Shortner-V1/core/config"
)

func GenerateMicroServiceAuthToken() (string, error) {
	configInstance := config.GetConfig()
	secretKey := []byte(configInstance.JwtMicroServiceSecretKey)

	claims := jwt.MapClaims{
		"sub":  "core",
		"type": "core-service",
		"exp":  timePackage.Now().Add(2 * timePackage.Minute).Unix(),
		"iat":  timePackage.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func GenerateAccessToken(userID string) (string, error) {
	configInstance := config.GetConfig()
	secretKey := []byte(configInstance.JwtSecretKey)

	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "access",
		"exp":  timePackage.Now().Add(15 * timePackage.Minute).Unix(),
		"iat":  timePackage.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userID string) (string, error) {
	configInstance := config.GetConfig()
	secretKey := []byte(configInstance.JwtRefreshKey)
	claims := jwt.MapClaims{
		"sub":  userID,
		"type": "refresh",
		"exp":  timePackage.Now().Add(7 * 24 * timePackage.Hour).Unix(),
		"iat":  timePackage.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func VerifyAccessToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	configInstance := config.GetConfig()
	secretKey := []byte(configInstance.JwtSecretKey)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Prevent algorithm spoofing
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, jwt.ErrTokenInvalidClaims
	}

	// Optional: enforce type
	if claims["type"] != "access" {
		return nil, nil, jwt.ErrTokenInvalidClaims
	}

	return token, claims, nil
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	configInstance := config.GetConfig()
	secretKey := []byte(configInstance.JwtRefreshKey)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, nil, jwt.ErrTokenInvalidClaims
	}

	if claims["type"] != "refresh" {
		return nil, nil, jwt.ErrTokenInvalidClaims
	}

	return token, claims, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	_, claims, err := VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	userID := claims["sub"].(string)

	return GenerateAccessToken(userID)
}
