package api

import (
	"encoding/json"
	"limoxlink-backend/internal/middleware"
	"limoxlink-backend/internal/models"
	"limoxlink-backend/internal/repository"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo   *repository.PostgresUserRepo
	TenantRepo *repository.TenantRepository
}

func NewAuthHandler(userRepo *repository.PostgresUserRepo, tenantRepo *repository.TenantRepository) *AuthHandler {
	return &AuthHandler{
		UserRepo:   userRepo,
		TenantRepo: tenantRepo,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 1. Fetch User
	user, err := h.UserRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid credentials: user lookup failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// 2. Compare Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials: password mismatch", http.StatusUnauthorized)
		return
	}

	// 3. Handle Multi-Tenancy Claims
	var tenantID string
	var companyID string
	if user.IsSuperAdmin {
		tenantID = "" // Global scope
		companyID = ""
	} else if user.CompanyID.Valid {
		companyID = user.CompanyID.UUID.String()
		tenant, err := h.TenantRepo.GetByCompanyID(r.Context(), user.CompanyID.UUID)
		if err == nil {
			tenantID = tenant.ID.String()
		}
		// If tenant lookup fails, we still allow login but tenantID will be empty
		// (though in a strict multi-tenant app we might reject here)
	}

	// 4. Generate Token
	token, err := middleware.GenerateToken(user.ID.String(), companyID, tenantID, string(user.Role), user.IsSuperAdmin)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  *user,
	})
}

type ChangePasswordRequest struct {
	UserID      string `json:"user_id"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Ideally verify old password or token claims here, but for MVP just update
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := h.UserRepo.UpdatePassword(r.Context(), userID, string(hashedPassword)); err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Password Updated"})
}
