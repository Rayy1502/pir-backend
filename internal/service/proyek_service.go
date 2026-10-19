package service

import (
	"context"

	"github.com/Rayy1502/pir-backend/internal/dto"
	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
)

type ProyekService struct {
	Repo *repository.ProyekRepository
}

func NewProyekService(repo *repository.ProyekRepository) *ProyekService {
	return &ProyekService{Repo: repo}
}

func (s *ProyekService) GetAll(ctx context.Context, f dto.ProyekFilter) ([]model.Proyek, error) {
	return s.Repo.GetAll(ctx, f)
}

func (s *ProyekService) GetByID(ctx context.Context, id string) (*model.Proyek, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ProyekService) GetDetailByID(ctx context.Context, id string) (*model.ProyekDetail, error) {
	return s.Repo.GetDetailByID(ctx, id)
}

func (s *ProyekService) GetStatistik(ctx context.Context) (*model.StatistikPIR, error) {
	return s.Repo.GetStatistik(ctx)
}

func (s *ProyekService) GetDashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	stat, err := s.Repo.GetStatistik(ctx)
	if err != nil {
		return nil, err
	}

	return &model.DashboardStats{
		TotalProyek:   stat.Total,
		ProyekSelesai: stat.Selesai,
		PIRPending:    stat.Pending + stat.Terlambat,
		PIRSelesai:    stat.Selesai,
	}, nil
}
