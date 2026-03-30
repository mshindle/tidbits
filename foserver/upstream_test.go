package foserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type mockUserService struct {
	findFunc func(id string) (*User, error)
}

func (m *mockUserService) FindByID(id string) (*User, error) {
	return m.findFunc(id)
}

func createToken(userID string, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func TestUserMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		mockSvc        *mockUserService
		expectedStatus int
		expectedBody   string
		checkContext   bool
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "missing authorization header",
		},
		{
			name:           "invalid authorization header format - no Bearer",
			authHeader:     "Basic dXNlcjpwYXNz",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid authorization header format",
		},
		{
			name:           "invalid authorization header format - empty token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid authorization header format",
		},
		{
			name:           "invalid token - bad secret",
			authHeader:     "Bearer " + createToken("123", "wrong-secret"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid token",
		},
		{
			name:           "invalid token - malformed",
			authHeader:     "Bearer malformed.token.here",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid token",
		},
		{
			name: "missing user_id claim",
			authHeader: "Bearer " + func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
				s, _ := token.SignedString([]byte(jwtSecret))
				return s
			}(),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "user_id not found in token",
		},
		{
			name:       "user not found - explicit error",
			authHeader: "Bearer " + createToken("404", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return nil, errors.New("user not found")
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name:       "user not found - wrapped error",
			authHeader: "Bearer " + createToken("404", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return nil, fmt.Errorf("lookup failed: %w", errors.New("User Not Found"))
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name:       "user not found - nil user",
			authHeader: "Bearer " + createToken("404", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return nil, nil
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name:       "internal service error",
			authHeader: "Bearer " + createToken("500", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return nil, errors.New("database connection lost")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal error",
		},
		{
			name:       "successful user retrieval",
			authHeader: "Bearer " + createToken("1", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return &User{ID: "1", Name: "John Doe", Email: "john@example.com"}, nil
				},
			},
			expectedStatus: http.StatusOK,
			checkContext:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkContext {
					user, ok := getUserFromContext(r.Context())
					if !ok || user == nil {
						t.Error("user not found in context")
					} else if user.ID != "1" {
						t.Errorf("expected user ID 1, got %s", user.ID)
					}
				}
				w.WriteHeader(http.StatusOK)
			})

			middleware := UserMiddleware(nextHandler, tt.mockSvc)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedBody != "" {
				if !strings.Contains(rec.Body.String(), tt.expectedBody) {
					t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rec.Body.String())
				}
			}
		})
	}
}

func TestGetUserMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		mockSvc        *mockUserService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "missing authorization header",
		},
		{
			name:       "any error becomes 404",
			authHeader: "Bearer " + createToken("500", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return nil, errors.New("database connection lost")
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found",
		},
		{
			name:       "successful user retrieval",
			authHeader: "Bearer " + createToken("1", jwtSecret),
			mockSvc: &mockUserService{
				findFunc: func(id string) (*User, error) {
					return &User{ID: "1", Name: "John Doe", Email: "john@example.com"}, nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := getUserMiddleware(nextHandler, tt.mockSvc)

			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			middleware.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedBody != "" {
				if !strings.Contains(rec.Body.String(), tt.expectedBody) {
					t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rec.Body.String())
				}
			}
		})
	}
}
