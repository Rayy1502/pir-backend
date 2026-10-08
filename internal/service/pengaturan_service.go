package service

import (
	"context"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
)

type PengaturanService struct {
	Repo *repository.PengaturanRepository
}

func NewPengaturanService(repo *repository.PengaturanRepository) *PengaturanService {
	return &PengaturanService{Repo: repo}
}

func (s *PengaturanService) GetAll(ctx context.Context) ([]model.PengaturanUnitKerja, error) {
	return s.Repo.GetAll(ctx)
}

func (s *PengaturanService) GetByUnitKerja(ctx context.Context, unit string) (*model.PengaturanUnitKerja, error) {
	return s.Repo.GetByUnitKerja(ctx, unit)
}

func (s *PengaturanService) Update(ctx context.Context, unit string, batasBulan int) error {
	if batasBulan < 1 || batasBulan > 24 {
		return errors.New("batas waktu harus antara 1-24 bulan")
	}

	existing, err := s.Repo.GetByUnitKerja(ctx, unit)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("unit kerja tidak ditemukan")
	}

	return s.Repo.Update(ctx, unit, batasBulan)
}
