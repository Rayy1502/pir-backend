package handler

import (
	"net/http"

	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type DashboardHandler struct {
	Service *service.ProyekService
}

func NewDashboardHandler(s *service.ProyekService) *DashboardHandler {
	return &DashboardHandler{Service: s}
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.Service.GetDashboardStats(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.Success(w, stats, nil)
}

func (h *DashboardHandler) StatistikPIR(w http.ResponseWriter, r *http.Request) {
	stat, err := h.Service.GetStatistik(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.Success(w, stat, nil)
}
