package middlewares

import (
	"bank-system/src/utils"
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type jwtMiddleware struct {
	logger *logrus.Logger
}

func NewJwtMiddleware(logger *logrus.Logger) *jwtMiddleware {
	return &jwtMiddleware{
		logger: logger,
	}
}

func (this *jwtMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			this.logger.Error("Authorization header is missing")
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		const bearerSchema = "Bearer "

		if !strings.HasPrefix(authHeader, bearerSchema) {
			this.logger.Error("Authorization header is invalid")
			http.Error(w, "Authorization header is invalid", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerSchema)

		userId, err := utils.ValidateJwt(tokenString)

		if err != nil {
			this.logger.Errorf("Failed to validate jwt: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
