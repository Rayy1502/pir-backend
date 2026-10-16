package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type PengaturanHandler struct {
	Service *service.PengaturanService
}

func NewPengaturanHandler(s *service.PengaturanService) *PengaturanHandler {
	return &PengaturanHandler{Service: s}
}

func (h *PengaturanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Service.GetAll(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.Success(w, list, nil)
}

func (h *PengaturanHandler) Update(w http.ResponseWriter, r *http.Request) {
	unit := chi.URLParam(r, "unit")
	if unescaped, err := url.PathUnescape(unit); err == nil && unescaped != "" {
		unit = unescaped
	}

	var body struct {
		BatasWaktuBulan int `json:"batas_waktu_bulan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	if err := h.Service.Update(r.Context(), unit, body.BatasWaktuBulan); err != nil {
		response.Error(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{"message": "Pengaturan berhasil disimpan"}, nil)
}
