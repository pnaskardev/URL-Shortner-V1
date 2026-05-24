package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/pnaskardev/URL-Shortner-V1/shortener-service/helpers/utils"
)

func MicroServiceAuthMiddleware(next http.Handler) http.Handler {
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
			fmt.Println(err)
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		slog.Debug("DEBUG", "CLAIMS", claims)

		// userID := claims["sub"]

		// ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
