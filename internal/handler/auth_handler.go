package handler

import (
	"clothing-store/internal/model"
	"clothing-store/internal/service"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" example:"user@example.com"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" example:"user@example.com"`
	Token       string `json:"token" example:"RESET123"`
	NewPassword string `json:"new_password" example:"newpassword123"`
}

type AuthHandler struct {
	AuthService *service.AuthService
}

// Login godoc
// @Summary      Login user
// @Description  Login with email and password to get JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login Credentials"
// @Success      200  {object}  LoginResponse
// @Failure      400  {string}  string "Invalid request"
// @Failure      401  {string}  string "Unauthorized"
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Register godoc
// @Summary      Register user
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration Info"
// @Success      201  {string}  string "User registered successfully"
// @Failure      400  {string}  string "Invalid request"
// @Failure      500  {string}  string "Internal server error"
// @Router       /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.AuthService.Register(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Request a password reset token via email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ForgotPasswordRequest true "Email address"
// @Success      200  {object}  map[string]string "message, token (for demo)"
// @Failure      400  {string}  string "Invalid request"
// @Router       /forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.ForgotPasswordRequest(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If the email exists, a reset token has been generated.",
		"token":   token, // Include token for demo purposes
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset password using the provided token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ResetPasswordRequest true "Reset details"
// @Success      200  {string}  string "Password reset successfully"
// @Failure      400  {string}  string "Invalid token or request"
// @Router       /reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.AuthService.ResetPassword(req.Email, req.Token, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Password reset successfully"})
}
