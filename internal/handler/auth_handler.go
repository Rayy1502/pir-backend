package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Rayy1502/pir-backend/internal/dto"
	"github.com/Rayy1502/pir-backend/internal/middleware"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body request tidak valid", nil)
		return
	}

	if req.NIP == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "NIP dan password wajib diisi", nil)
		return
	}

	user, accessToken, refreshToken, err := h.AuthService.Login(r.Context(), req.NIP, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error(), nil)
		return
	}

	response.Success(w, dto.LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User tidak terautentikasi", nil)
		return
	}

	user, err := h.AuthService.GetUserByID(r.Context(), userID)
	if err != nil || user == nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "User tidak ditemukan", nil)
		return
	}

	response.Success(w, user, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Untuk v1.0: stateless logout (client hapus token)
	response.Success(w, map[string]string{"message": "Berhasil logout"}, nil)
}
