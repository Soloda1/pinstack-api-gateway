package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"pinstack-api-gateway/internal/custom_errors"
	"pinstack-api-gateway/internal/logger"
	"pinstack-api-gateway/internal/utils"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

func JWTValidationMiddleware(secretKey string, log *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			entry := log.With(
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				entry.Error("authorization header is missing")
				utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				entry.Error("invalid authorization header format")
				utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidToken.Error())
				return
			}

			parser := jwt.NewParser(
				jwt.WithValidMethods([]string{"HS256"}),
				jwt.WithLeeway(5*time.Second),
			)

			token, err := parser.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil {
				entry.Error("token validation failed", slog.String("error", err.Error()))
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrTokenExpired.Error())
				case errors.Is(err, jwt.ErrTokenMalformed):
					utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidToken.Error())
				default:
					utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrUnauthenticated.Error())
				}
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				entry.Error("invalid token claims")
				utils.SendError(w, http.StatusUnauthorized, custom_errors.ErrInvalidToken.Error())
				return
			}

			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))

			entry.Info("token validation completed successfully", slog.Int64("user_id", claims.UserID))
		}

		return http.HandlerFunc(fn)
	}
}

func GetClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok {
		return nil, custom_errors.ErrUnauthenticated
	}
	return claims, nil
}
