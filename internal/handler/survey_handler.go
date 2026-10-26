package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type SurveyHandler struct {
	Service *service.SurveyService
}

func NewSurveyHandler(s *service.SurveyService) *SurveyHandler {
	return &SurveyHandler{Service: s}
}

func (h *SurveyHandler) Save(w http.ResponseWriter, r *http.Request) {
	pirID := chi.URLParam(r, "id")
	if pirID == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "PIR ID dibutuhkan", nil)
		return
	}

	var req model.SubmitSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	if err := h.Service.SaveSurvey(r.Context(), pirID, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "SAVE_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{
		"message": "Survey berhasil disimpan",
	}, nil)
}

func (h *SurveyHandler) Get(w http.ResponseWriter, r *http.Request) {
	pirID := chi.URLParam(r, "id")
	if pirID == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "PIR ID dibutuhkan", nil)
		return
	}

	survey, err := h.Service.GetSurvey(r.Context(), pirID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, survey, nil)
}
