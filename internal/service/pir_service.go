package service

import (
	"context"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
)

type PIRService struct {
	Repo *repository.PIRRepository
}

func NewPIRService(repo *repository.PIRRepository) *PIRService {
	return &PIRService{Repo: repo}
}

func (s *PIRService) Create(ctx context.Context, req *model.CreatePIRRequest, createdBy string) (*model.PIRForm, error) {
	if req.ProyekID == "" {
		return nil, errors.New("proyek_id wajib diisi")
	}
	if req.Status == "submitted" {
		if req.Vendor == "" {
			return nil, errors.New("vendor/tim pengembang wajib diisi")
		}
		if req.TujuanAwal == "" {
			return nil, errors.New("tujuan awal pengembangan wajib diisi")
		}
		if req.Metode == "" {
			return nil, errors.New("metode evaluasi PIR wajib dipilih")
		}
	}
	return s.Repo.Create(ctx, req, createdBy)
}

func (s *PIRService) GetByID(ctx context.Context, id string) (*model.PIRForm, error) {
	if id == "" {
		return nil, errors.New("id tidak boleh kosong")
	}
	form, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form == nil {
		return nil, errors.New("form PIR tidak ditemukan")
	}
	return form, nil
}

func (s *PIRService) GetByProyekID(ctx context.Context, proyekID string) (*model.PIRForm, error) {
	if proyekID == "" {
		return nil, errors.New("proyek_id tidak boleh kosong")
	}
	return s.Repo.GetByProyekID(ctx, proyekID)
}

func (s *PIRService) Update(ctx context.Context, id string, req *model.UpdatePIRRequest) error {
	if id == "" {
		return errors.New("id tidak boleh kosong")
	}
	if req.Status == "submitted" {
		if req.Vendor == "" {
			return errors.New("vendor/tim pengembang wajib diisi")
		}
		if req.TujuanAwal == "" {
			return errors.New("tujuan awal pengembangan wajib diisi")
		}
		if req.Metode == "" {
			return errors.New("metode evaluasi PIR wajib dipilih")
		}
	}
	return s.Repo.Update(ctx, id, req)
}

func (s *PIRService) Approve(ctx context.Context, id, approvedBy string) error {
	if id == "" {
		return errors.New("id tidak boleh kosong")
	}
	return s.Repo.Approve(ctx, id, approvedBy)
}
