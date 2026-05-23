package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/pnaskardev/URL-Shortner-V1/core/helpers/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// SPlit at space
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		_, claims, err := utils.VerifyAccessToken(token)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		userID := claims["sub"]

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
