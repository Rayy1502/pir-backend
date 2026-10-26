package service

import (
	"context"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
	"github.com/Rayy1502/pir-backend/pkg/pdf"
)

type LaporanService struct {
	laporanRepo *repository.LaporanRepository
	proyekRepo  *repository.ProyekRepository
}

func NewLaporanService(laporanRepo *repository.LaporanRepository, proyekRepo *repository.ProyekRepository) *LaporanService {
	return &LaporanService{
		laporanRepo: laporanRepo,
		proyekRepo:  proyekRepo,
	}
}

func (s *LaporanService) Create(ctx context.Context, req *model.CreateLaporanRequest, pembuatID string) (*model.LaporanPIR, error) {
	if req.ProyekID == "" {
		return nil, errors.New("proyek_id is required")
	}
	if req.PIRFormID == "" {
		return nil, errors.New("pir_form_id is required")
	}

	return s.laporanRepo.Create(ctx, req, pembuatID)
}

func (s *LaporanService) GetByID(ctx context.Context, id string) (*model.LaporanPIR, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.laporanRepo.GetByID(ctx, id)
}

func (s *LaporanService) GetByProyekID(ctx context.Context, proyekID string) (*model.LaporanPIR, error) {
	if proyekID == "" {
		return nil, errors.New("proyek_id is required")
	}
	return s.laporanRepo.GetByProyekID(ctx, proyekID)
}

func (s *LaporanService) GetAll(ctx context.Context, proyekID, search string) ([]model.LaporanPIR, error) {
	return s.laporanRepo.GetAll(ctx, proyekID, search)
}

func (s *LaporanService) GeneratePDF(ctx context.Context, id string) ([]byte, *model.LaporanPIR, error) {
	lap, err := s.laporanRepo.GetByID(ctx, id)
	if err != nil || lap == nil {
		return nil, nil, errors.New("laporan not found")
	}

	pdfBytes, err := pdf.GenerateLaporanPDF(lap)
	if err != nil {
		return nil, nil, err
	}

	return pdfBytes, lap, nil
}

func (s *LaporanService) SaveSignedPDF(ctx context.Context, id, fileURL, fileName string) (*model.LaporanPIR, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if fileURL == "" {
		return nil, errors.New("fileURL is required")
	}
	return s.laporanRepo.SaveSignedPDF(ctx, id, fileURL, fileName)
}
