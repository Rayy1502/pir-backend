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

type PIRHandler struct {
	Service   *service.PIRService
	UploadDir string
}

func NewPIRHandler(s *service.PIRService, uploadDir string) *PIRHandler {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	return &PIRHandler{Service: s, UploadDir: uploadDir}
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

func (h *PIRHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "ID PIR form dibutuhkan", nil)
		return
	}

	var req model.ApprovePIRRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.ApprovedBy == "" {
		req.ApprovedBy = "VP Teknologi Informasi"
	}

	if err := h.Service.Approve(r.Context(), id, req.ApprovedBy); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{
		"message":     "Data evaluasi PIR berhasil disetujui (ACC)",
		"approved_by": req.ApprovedBy,
	}, nil)
}

func (h *PIRHandler) GetDokumentasi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "ID form PIR dibutuhkan", nil)
		return
	}
	kategori := r.URL.Query().Get("kategori")
	list, err := h.Service.GetDokumentasi(r.Context(), id, kategori)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}
	response.Success(w, list, nil)
}

func (h *PIRHandler) UploadDokumentasi(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "ID form PIR dibutuhkan", nil)
		return
	}

	// Batasi ukuran request max 15MB
	r.Body = http.MaxBytesReader(w, r.Body, 15<<20)
	if err := r.ParseMultipartForm(15 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_TOO_LARGE", "Ukuran file melebihi batas maksimal 15MB", nil)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "NO_FILE", "File dokumentasi belum dipilih", nil)
		return
	}
	defer file.Close()

	kategori := r.FormValue("kategori")
	if kategori == "" {
		kategori = "observasi"
	}
	keterangan := r.FormValue("keterangan")

	uploadDir := filepath.Join(h.UploadDir, "dokumentasi")
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

	size, err := io.Copy(dst, file)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "IO_ERROR", "Gagal menyalin isi file", nil)
		return
	}

	fileURL := fmt.Sprintf("/uploads/dokumentasi/%s", storedName)

	dok, err := h.Service.AddDokumentasi(r.Context(), id, kategori, origName, fileURL, size, keterangan)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, dok, nil)
}

func (h *PIRHandler) DeleteDokumentasi(w http.ResponseWriter, r *http.Request) {
	dokID := chi.URLParam(r, "dokId")
	if dokID == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "ID dokumentasi dibutuhkan", nil)
		return
	}

	if err := h.Service.DeleteDokumentasi(r.Context(), dokID); err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, map[string]string{"message": "Dokumentasi berhasil dihapus"}, nil)
}

