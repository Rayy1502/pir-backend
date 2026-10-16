package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Rayy1502/pir-backend/internal/middleware"
	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type PIRHandler struct {
	Service *service.PIRService
}

func NewPIRHandler(s *service.PIRService) *PIRHandler {
	return &PIRHandler{Service: s}
}

func (h *PIRHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePIRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	userID := middleware.GetUserID(r.Context())
	form, err := h.Service.Create(r.Context(), &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, form, nil)
}

func (h *PIRHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	form, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	response.Success(w, form, nil)
}

func (h *PIRHandler) GetByProyekOrList(w http.ResponseWriter, r *http.Request) {
	proyekID := r.URL.Query().Get("proyek_id")
	if proyekID != "" {
		form, err := h.Service.GetByProyekID(r.Context(), proyekID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
			return
		}
		if form == nil {
			response.Success(w, nil, nil)
			return
		}
		response.Success(w, form, nil)
		return
	}

	response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Query proyek_id dibutuhkan", nil)
}

func (h *PIRHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req model.UpdatePIRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	if err := h.Service.Update(r.Context(), id, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{"message": "Form PIR berhasil diperbarui"}, nil)
}
