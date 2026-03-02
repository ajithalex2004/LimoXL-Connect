package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const ClaimsKey = "user"

var SecretKey = []byte(getSecret())

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "super-secret-key-change-this-in-prod"
	}
	return secret
}

type Claims struct {
	UserID       string `json:"user_id"`
	CompanyID    string `json:"company_id"`
	TenantID     string `json:"tenant_id"`
	Role         string `json:"role"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID, companyID, tenantID, role string, isSuperAdmin bool) (string, error) {
	claims := &Claims{
		UserID:       userID,
		CompanyID:    companyID,
		TenantID:     tenantID,
		Role:         role,
		IsSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// SuperAdminMiddleware restricts access to Super Admins only
func SuperAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*Claims)
		if !ok || !claims.IsSuperAdmin {
			http.Error(w, "SuperAdmin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ITenantRepo is an interface for FeatureGate to avoid circular dependency
type ITenantRepo interface {
	IsFeatureEnabled(ctx context.Context, tenantID uuid.UUID, featureKey string) (bool, error)
}

// FeatureGate restricting access to specific microservices based on tenant settings
func FeatureGate(repo ITenantRepo, featureKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// SuperAdmins bypass feature gates
			if claims.IsSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}

			if claims.TenantID == "" {
				http.Error(w, "Tenant identity required", http.StatusForbidden)
				return
			}

			tenantID, err := uuid.Parse(claims.TenantID)
			if err != nil {
				http.Error(w, "Invalid tenant ID", http.StatusForbidden)
				return
			}

			enabled, err := repo.IsFeatureEnabled(r.Context(), tenantID, featureKey)
			if err != nil || !enabled {
				http.Error(w, fmt.Sprintf("Feature '%s' is not enabled for your plan", featureKey), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware validates the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to context directly
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
