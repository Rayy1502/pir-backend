package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Rayy1502/pir-backend/internal/dto"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type ProyekHandler struct {
	Service *service.ProyekService
}

func NewProyekHandler(s *service.ProyekService) *ProyekHandler {
	return &ProyekHandler{Service: s}
}

func (h *ProyekHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	f := dto.ProyekFilter{
		Query:    r.URL.Query().Get("q"),
		Kategori: r.URL.Query().Get("kategori"),
		Status:   r.URL.Query().Get("status"),
	}

	list, err := h.Service.GetAll(r.Context(), f)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.Success(w, list, nil)
}

func (h *ProyekHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := h.Service.GetDetailByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	if p == nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Proyek tidak ditemukan", nil)
		return
	}
	response.Success(w, p, nil)
}

func (h *ProyekHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProyekRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_BODY", "Request body tidak valid", nil)
		return
	}

	if req.NamaKegiatan == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Nama kegiatan wajib diisi", nil)
		return
	}
	if req.Kategori == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Kategori proyek wajib diisi", nil)
		return
	}
	if req.UnitKerja == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Unit kerja pemilik wajib diisi", nil)
		return
	}
	if req.TanggalGoLive == "" {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Tanggal go live wajib diisi", nil)
		return
	}

	p, err := h.Service.Create(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Created(w, p)
}


