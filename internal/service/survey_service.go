package service

import (
	"context"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
)

type SurveyService struct {
	surveyRepo *repository.SurveyRepository
	pirRepo    *repository.PIRRepository
}

func NewSurveyService(surveyRepo *repository.SurveyRepository, pirRepo *repository.PIRRepository) *SurveyService {
	return &SurveyService{
		surveyRepo: surveyRepo,
		pirRepo:    pirRepo,
	}
}

func (s *SurveyService) SaveSurvey(ctx context.Context, pirFormID string, req *model.SubmitSurveyRequest) error {
	if pirFormID == "" {
		return errors.New("pir_form_id is required")
	}

	pir, err := s.pirRepo.GetByID(ctx, pirFormID)
	if err != nil || pir == nil {
		return errors.New("pir_form_not_found")
	}

	if req.Kepuasan < 1 || req.Kepuasan > 5 {
		req.Kepuasan = 4
	}
	if req.KesesuaianTujuan < 1 || req.KesesuaianTujuan > 5 {
		req.KesesuaianTujuan = 4
	}

	return s.surveyRepo.SaveSurvey(ctx, pirFormID, req)
}

func (s *SurveyService) GetSurvey(ctx context.Context, pirFormID string) (*model.SurveyResponse, error) {
	if pirFormID == "" {
		return nil, errors.New("pir_form_id is required")
	}

	res, err := s.surveyRepo.GetByPIRFormID(ctx, pirFormID)
	if err != nil {
		if err.Error() == "survey_not_found" {
			// Return default empty template so client can fill it immediately
			return &model.SurveyResponse{
				PIRFormID:        pirFormID,
				Kepuasan:         5,
				KesesuaianTujuan: 5,
				Kendala:          "",
				Saran:            "",
				Pertanyaan:       []model.SurveyPertanyaan{},
			}, nil
		}
		return nil, err
	}

	return res, nil
}
