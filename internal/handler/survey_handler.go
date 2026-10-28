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

// SaveQuestions menyimpan daftar pertanyaan yang dikonfigurasi PIC
func (h *SurveyHandler) SaveQuestions(w http.ResponseWriter, r *http.Request) {
	pirID := chi.URLParam(r, "id")
	if pirID == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "PIR ID dibutuhkan", nil)
		return
	}

	var req model.SaveQuestionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid: "+err.Error(), nil)
		return
	}

	if err := h.Service.SaveQuestions(r.Context(), pirID, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "SAVE_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{
		"message": "Daftar pertanyaan survey berhasil disimpan",
	}, nil)
}

// GetPublic mengambil form survey publik untuk user / responden
func (h *SurveyHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	pirID := chi.URLParam(r, "id")
	if pirID == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "PIR ID dibutuhkan", nil)
		return
	}

	survey, err := h.Service.GetPublicSurvey(r.Context(), pirID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Survey tidak ditemukan atau link sudah kedaluwarsa", nil)
		return
	}

	response.Success(w, survey, nil)
}

// PublicSubmit menangani pengiriman review oleh responden dari link publik
func (h *SurveyHandler) PublicSubmit(w http.ResponseWriter, r *http.Request) {
	pirID := chi.URLParam(r, "id")
	if pirID == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "PIR ID dibutuhkan", nil)
		return
	}

	var req model.PublicSubmitSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Format data kuesioner tidak valid", nil)
		return
	}

	if err := h.Service.PublicSubmit(r.Context(), pirID, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "SUBMIT_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{
		"message": "Review survey berhasil dikirim. Terima kasih atas partisipasi Anda!",
	}, nil)
}

// Save menangani pengisian manual jika PIC memilih mengisi sendiri (fallback)
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

// Get mengambil data survey lengkap untuk PIC di dashboard
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
