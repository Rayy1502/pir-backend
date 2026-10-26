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

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/service"
	"github.com/Rayy1502/pir-backend/pkg/response"
)

type TindakLanjutHandler struct {
	Service   *service.TindakLanjutService
	UploadDir string
}

func NewTindakLanjutHandler(s *service.TindakLanjutService, uploadDir string) *TindakLanjutHandler {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	_ = os.MkdirAll(uploadDir, 0755)
	return &TindakLanjutHandler{
		Service:   s,
		UploadDir: uploadDir,
	}
}

func (h *TindakLanjutHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	laporanID := r.URL.Query().Get("laporan_id")
	search := r.URL.Query().Get("q")

	list, err := h.Service.GetAll(r.Context(), status, laporanID, search)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, list, nil)
}

func (h *TindakLanjutHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, err := h.Service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", err.Error(), nil)
		return
	}
	if item == nil {
		response.Error(w, http.StatusNotFound, "NOT_FOUND", "Tindak lanjut tidak ditemukan", nil)
		return
	}

	response.Success(w, item, nil)
}

func (h *TindakLanjutHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTindakLanjutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	item, err := h.Service.Create(r.Context(), &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "CREATE_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, item, nil)
}

func (h *TindakLanjutHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req model.UpdateTindakLanjutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "Body tidak valid", nil)
		return
	}

	item, err := h.Service.Update(r.Context(), id, &req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "UPDATE_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, item, nil)
}

func (h *TindakLanjutHandler) UploadBukti(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "BAD_REQUEST", "ID tindak lanjut dibutuhkan", nil)
		return
	}

	// Batasi ukuran request max 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "FILE_TOO_LARGE", "Ukuran file melebihi batas maksimal 10MB", nil)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "NO_FILE", "File bukti belum diupload", nil)
		return
	}
	defer file.Close()

	_ = os.MkdirAll(h.UploadDir, 0755)

	origName := filepath.Base(fileHeader.Filename)
	safeName := strings.ReplaceAll(origName, " ", "_")
	storedName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), safeName)
	targetPath := filepath.Join(h.UploadDir, storedName)

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

	fileURL := fmt.Sprintf("/uploads/%s", storedName)

	bukti, err := h.Service.AddBukti(r.Context(), id, origName, fileURL, size)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "DB_ERROR", err.Error(), nil)
		return
	}

	response.Success(w, bukti, nil)
}
