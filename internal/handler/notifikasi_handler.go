package handler

import (
	"net/http"
	"strconv"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
	"github.com/go-chi/chi/v5"
)

type NotifikasiHandler struct {
	service *service.NotifikasiService
}

func NewNotifikasiHandler(service *service.NotifikasiService) *NotifikasiHandler {
	return &NotifikasiHandler{service: service}
}

// GetAll mengambil daftar notifikasi
func (h *NotifikasiHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var filter model.NotifikasiFilter

	dibacaParam := r.URL.Query().Get("dibaca")
	if dibacaParam != "" {
		d, err := strconv.ParseBool(dibacaParam)
		if err == nil {
			filter.Dibaca = &d
		}
	}

	list, err := h.service.GetAll(r.Context(), filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_ERROR", "Gagal mengambil daftar notifikasi: "+err.Error(), nil)
		return
	}

	response.Success(w, list, nil)
}

// GetUnreadCount mengambil jumlah notifikasi belum dibaca
func (h *NotifikasiHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.GetUnreadCount(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "FETCH_ERROR", "Gagal mengambil jumlah notifikasi", nil)
		return
	}

	response.Success(w, map[string]int{"unread_count": count}, nil)
}

// MarkAsRead menandai satu notifikasi sebagai telah dibaca
func (h *NotifikasiHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "INVALID_ID", "ID notifikasi wajib diisi", nil)
		return
	}

	if err := h.service.MarkAsRead(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "UPDATE_ERROR", "Gagal menandai notifikasi telah dibaca", nil)
		return
	}

	response.Success(w, map[string]string{"message": "Notifikasi ditandai telah dibaca"}, nil)
}

// MarkAllAsRead menandai seluruh notifikasi sebagai telah dibaca
func (h *NotifikasiHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	if err := h.service.MarkAllAsRead(r.Context()); err != nil {
		response.Error(w, http.StatusInternalServerError, "UPDATE_ERROR", "Gagal menandai seluruh notifikasi telah dibaca", nil)
		return
	}

	response.Success(w, map[string]string{"message": "Semua notifikasi ditandai telah dibaca"}, nil)
}

// TriggerEvaluation memicu evaluasi notifikasi secara langsung
func (h *NotifikasiHandler) TriggerEvaluation(w http.ResponseWriter, r *http.Request) {
	if err := h.service.TriggerEvaluation(r.Context()); err != nil {
		response.Error(w, http.StatusInternalServerError, "SCHEDULER_ERROR", "Gagal memicu evaluasi notifikasi: "+err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{"message": "Evaluasi notifikasi background berhasil dipicu"}, nil)
}
