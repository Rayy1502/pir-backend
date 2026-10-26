package service

import (
	"context"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
)

type TindakLanjutService struct {
	repo *repository.TindakLanjutRepository
}

func NewTindakLanjutService(repo *repository.TindakLanjutRepository) *TindakLanjutService {
	return &TindakLanjutService{repo: repo}
}

func (s *TindakLanjutService) GetAll(ctx context.Context, status, laporanID, search string) ([]model.TindakLanjut, error) {
	return s.repo.GetAll(ctx, status, laporanID, search)
}

func (s *TindakLanjutService) GetByID(ctx context.Context, id string) (*model.TindakLanjut, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TindakLanjutService) Create(ctx context.Context, req *model.CreateTindakLanjutRequest) (*model.TindakLanjut, error) {
	if req.LaporanID == "" {
		return nil, errors.New("laporan_id is required")
	}
	if req.Temuan == "" {
		return nil, errors.New("temuan is required")
	}
	return s.repo.Create(ctx, req)
}

func (s *TindakLanjutService) Update(ctx context.Context, id string, req *model.UpdateTindakLanjutRequest) (*model.TindakLanjut, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if err := s.repo.Update(ctx, id, req); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TindakLanjutService) AddBukti(ctx context.Context, tindakLanjutID, namaFile, url string, ukuranBytes int64) (*model.TindakLanjutBukti, error) {
	if tindakLanjutID == "" {
		return nil, errors.New("tindak_lanjut_id is required")
	}
	return s.repo.AddBukti(ctx, tindakLanjutID, namaFile, url, ukuranBytes)
}
