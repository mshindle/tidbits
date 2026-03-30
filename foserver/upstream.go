package foserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/golang-jwt/jwt/v5"
)

// listenAndServe sets up a mini web server that serves a predetermined response.
func listenAndServe(addr string, statusCode int, node int) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			output := fmt.Sprintf("{\"node-%02d\": %d}", node, node)
			_, _ = w.Write([]byte(output))
		}
		return
	})
	server := &http.Server{Addr: addr, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).WithField("node", node).Error("upstream node failed")
		}
	}()

	return server
}

type User struct {
	ID    string
	Name  string
	Email string
}

type UserService interface {
	FindByID(id string) (*User, error)
}

type contextKey string

const jwtSecret = "your-secret-key"
const userContextKey contextKey = "user"

// getUserMiddleware is badly AI generated middleware function. It automatically maps all UserService errors to
// / http.StatusNotFound instead of potentially http.InternalServerError
func getUserMiddleware(next http.Handler, svc UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			http.Error(w, "user_id not found in token", http.StatusUnauthorized)
			return
		}

		user, err := svc.FindByID(userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// UserMiddleware extracts the bearer token, validates it using jwtSecret, loads the user via svc,
// and stores the user in the request context under userContextKey.
//
// It distinguishes "user does not exist" from "internal service error":
// - missing user: returns 404
// - other svc errors: returns 500
func UserMiddleware(next http.Handler, svc UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || token == nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || strings.TrimSpace(userID) == "" {
			http.Error(w, "user_id not found in token", http.StatusUnauthorized)
			return
		}

		user, err := svc.FindByID(userID)
		if err != nil {
			if isUserNotFoundErr(err) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			// Treat "no error but no user" as "does not exist".
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isUserNotFoundErr(err error) bool {
	if err == nil {
		return false
	}

	// Walk the unwrap chain and look for the message text.
	// (Service guarantees "user not found" appears in the error or wrapped error.)
	for e := err; e != nil; e = errors.Unwrap(e) {
		if strings.Contains(strings.ToLower(e.Error()), "user not found") {
			return true
		}
	}
	return false
}
