package service

import (
	"context"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
	"github.com/Rayy1502/pir-backend/internal/scheduler"
)

type NotifikasiService struct {
	repo      *repository.NotifikasiRepository
	scheduler *scheduler.NotificationScheduler
}

func NewNotifikasiService(repo *repository.NotifikasiRepository, sched *scheduler.NotificationScheduler) *NotifikasiService {
	return &NotifikasiService{
		repo:      repo,
		scheduler: sched,
	}
}

func (s *NotifikasiService) GetAll(ctx context.Context, filter model.NotifikasiFilter) ([]model.Notifikasi, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s *NotifikasiService) GetUnreadCount(ctx context.Context) (int, error) {
	return s.repo.GetUnreadCount(ctx)
}

func (s *NotifikasiService) MarkAsRead(ctx context.Context, id string) error {
	return s.repo.MarkAsRead(ctx, id)
}

func (s *NotifikasiService) MarkAllAsRead(ctx context.Context) error {
	return s.repo.MarkAllAsRead(ctx)
}

func (s *NotifikasiService) TriggerEvaluation(ctx context.Context) error {
	if s.scheduler != nil {
		return s.scheduler.Evaluate(ctx)
	}
	return nil
}
