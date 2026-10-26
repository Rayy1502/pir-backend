package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/Rayy1502/pir-backend/internal/middleware"
	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type LaporanHandler struct {
	Service   *service.LaporanService
	UploadDir string
}

func NewLaporanHandler(s *service.LaporanService, uploadDir string) *LaporanHandler {
	return &LaporanHandler{
		Service:   s,
		UploadDir: uploadDir,
	}
}

func (h *LaporanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLaporanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	userID := middleware.GetUserID(r.Context())
	laporan, err := h.Service.Create(r.Context(), &req, userID)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, laporan, nil)
}

func (h *LaporanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	laporan, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	if laporan == nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Laporan tidak ditemukan", nil)
		return
	}

	response.Success(w, laporan, nil)
}

func (h *LaporanHandler) GetByProyekID(w http.ResponseWriter, r *http.Request) {
	proyekID := chi.URLParam(r, "proyekId")
	laporan, err := h.Service.GetByProyekID(r.Context(), proyekID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	if laporan == nil {
		response.Success(w, nil, nil)
		return
	}

	response.Success(w, laporan, nil)
}

func (h *LaporanHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	proyekID := r.URL.Query().Get("proyek_id")
	search := r.URL.Query().Get("q")

	list, err := h.Service.GetAll(r.Context(), proyekID, search)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, list, nil)
}

func (h *LaporanHandler) ExportPDF(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pdfBytes, lap, err := h.Service.GeneratePDF(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "PDF_ERROR", err.Error(), nil)
		return
	}

	safeNo := strings.ReplaceAll(lap.NomorPIR, "/", "_")
	filename := fmt.Sprintf("Laporan_%s.pdf", safeNo)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdfBytes)
}

func (h *LaporanHandler) UploadSignedPDF(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "ID laporan dibutuhkan", nil)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 15<<20) // max 15MB
	if err := r.ParseMultipartForm(15 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_TOO_LARGE", "Ukuran file melebihi batas 15MB", nil)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "NO_FILE", "File PDF bertanda tangan belum diunggah", nil)
		return
	}
	defer file.Close()

	uploadDir := filepath.Join(h.UploadDir, "laporan_ttd")
	_ = os.MkdirAll(uploadDir, 0755)

	origName := filepath.Base(fileHeader.Filename)
	safeName := strings.ReplaceAll(origName, " ", "_")
	storedName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), safeName)
	targetPath := filepath.Join(uploadDir, storedName)

	dst, err := os.Create(targetPath)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "IO_ERROR", "Gagal menyimpan file", nil)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "IO_ERROR", "Gagal menyalin isi file", nil)
		return
	}

	fileURL := fmt.Sprintf("/uploads/laporan_ttd/%s", storedName)

	laporan, err := h.Service.SaveSignedPDF(r.Context(), id, fileURL, origName)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, laporan, nil)
}
