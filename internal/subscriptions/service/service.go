package service

import (
	"context"
	"ef_mob_test_go/internal/models"
	"ef_mob_test_go/internal/subscriptions/repository"

	"go.uber.org/zap"
)

type subService struct {
	subRepo repository.SubRepo
	logger  *zap.SugaredLogger
}

type SubService interface {
	CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
}

func NewSubService(subRepo repository.SubRepo, log *zap.SugaredLogger) SubService {
	return &subService{
		subRepo: subRepo,
		logger:  log,
	}
}

func (s *subService) CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	createdSub, err := s.subRepo.CreateSub(ctx, sub)
	if err != nil {
		return nil, err
	}
	return createdSub, nil
}
