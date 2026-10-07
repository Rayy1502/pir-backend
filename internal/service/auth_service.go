package service

import (
	"context"
	"errors"
	"time"

	"github.com/Rayy1502/pir-backend/internal/config"
	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
	"github.com/Rayy1502/pir-backend/pkg/hash"
	"github.com/Rayy1502/pir-backend/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("NIP atau password salah")
	ErrUserInactive       = errors.New("Akun tidak aktif")
)

type AuthService struct {
	Cfg      *config.Config
	UserRepo *repository.UserRepository
}

func NewAuthService(cfg *config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{Cfg: cfg, UserRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, nip, password string) (*model.User, string, string, error) {
	user, err := s.UserRepo.FindByNIP(ctx, nip)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", ErrInvalidCredentials
	}

	if !hash.CheckPassword(user.PasswordHash, password) {
		return nil, "", "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", "", ErrUserInactive
	}

	accessToken, err := jwt.GenerateToken(s.Cfg.JWTSecret, user.ID, user.NIP, user.Role, 15*time.Minute)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := jwt.GenerateToken(s.Cfg.JWTSecret, user.ID, user.NIP, user.Role, 7*24*time.Hour)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.UserRepo.FindByID(ctx, id)
}
