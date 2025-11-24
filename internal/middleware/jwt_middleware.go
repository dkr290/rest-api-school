package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dkr290/go-advanced-projects/rest-api-school-management/config"
	"github.com/dkr290/go-advanced-projects/rest-api-school-management/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

func JWTMiddleware(next http.Handler, conf config.Config, logger logging.Logger) http.Handler {
	logger.Logging.Debugln(strings.Repeat("-", 20) + "JWT Middleware" + strings.Repeat("-", 20))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path
		for _, path := range conf.ExcludedAuthMiddlewarePath {
			path = strings.TrimSpace(path)
			if strings.HasPrefix(requestPath, path) {
				next.ServeHTTP(w, r)
				return // Important: Return so we don't execute the JWT check below
			}
		}

		token, err := r.Cookie("Bearer")
		if err != nil {
			logger.Logging.Debugln("No Bearer tioken found")
			http.Error(w, "Authorization Header Missing ", http.StatusUnauthorized)
			return
		}
		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Logging.Debugln("unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(conf.JWTSecret), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {

				logger.Logging.Error("Token is expired")
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			} else if errors.Is(err, jwt.ErrTokenMalformed) {
				logger.Logging.Error("Token is not valid")
				http.Error(w, "Token is not valid", http.StatusUnauthorized)
				return

			}

			logger.Logging.Debugf("error from middleware %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if parsedToken.Valid {
			logger.Logging.Info("Valid JWT token")
		} else {
			logger.Logging.Error("Invalid JWT token")
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if ok {
			logger.Logging.Debug(
				claims["uid"],
				claims["exp"],
				claims["role"],
			)
		} else {
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, ContextKey("expiresAt"), claims["exp"])
		ctx = context.WithValue(ctx, ContextKey("username"), claims["user"])
		ctx = context.WithValue(ctx, ContextKey("uid"), claims["uid"])

		logger.Logging.Debug(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))

		fmt.Println("Responce from JWT Middleware")
	})
}
