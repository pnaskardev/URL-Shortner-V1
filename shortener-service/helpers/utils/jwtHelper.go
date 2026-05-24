package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/config"
)

func VerifyAccessToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	configInstance := config.GetConfig()
	secretKey := []byte(configInstance.JwtMicroServiceSecretKey)

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

	return token, claims, nil
}
