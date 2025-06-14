package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextRole   contextKey = "role"
)

func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims := jwt.MapClaims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
				return
			}

			role := claims["role"].(string)
			userID := claims["user_id"].(string)

			authorized := false
			for _, allowed := range allowedRoles {
				if role == allowed {
					authorized = true
					break
				}
			}

			if !authorized {
				http.Error(w, "403 Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, userID)
			ctx = context.WithValue(ctx, ContextRole, role)
			next(w, r.WithContext(ctx))
		}
	}
}
