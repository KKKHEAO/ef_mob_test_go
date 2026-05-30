package service

import (
	"context"
	"fmt"
	"time"

	"ef_mob_test_go/internal/models"
	"ef_mob_test_go/internal/subscriptions/repository"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type subService struct {
	subRepo repository.SubRepo
	logger  *zap.SugaredLogger
}

type SubService interface {
	CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	GetSubByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	UpdateSubByID(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	DeleteSubByID(ctx context.Context, id uuid.UUID) error
	ListSubs(ctx context.Context, page, pageSize int) ([]models.Subscription, int, error)
	CalculateCost(ctx context.Context, filter models.PeriodFilter) (*models.CostResponse, error)
}

func NewSubService(subRepo repository.SubRepo, log *zap.SugaredLogger) SubService {
	return &subService{
		subRepo: subRepo,
		logger:  log,
	}
}

func (s *subService) CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate uuid: %w", err)
	}
	sub.ID = id

	createdSub, err := s.subRepo.CreateSub(ctx, sub)
	if err != nil {
		s.logger.Errorf("CreateSub: %v", err)
		return nil, err
	}
	s.logger.Infow("subscription created", "id", createdSub.ID, "name", createdSub.Name)
	return createdSub, nil
}

func (s *subService) GetSubByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	sub, err := s.subRepo.GetSubByID(ctx, id)
	if err != nil {
		s.logger.Errorf("GetSubByID: %v", err)
		return nil, err
	}
	return sub, nil
}

func (s *subService) UpdateSubByID(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	updated, err := s.subRepo.UpdateSubByID(ctx, sub)
	if err != nil {
		s.logger.Errorf("UpdateSubByID: %v", err)
		return nil, err
	}
	s.logger.Infow("subscription updated", "id", updated.ID, "name", updated.Name)
	return updated, nil
}

func (s *subService) ListSubs(ctx context.Context, page, pageSize int) ([]models.Subscription, int, error) {
	offset := (page - 1) * pageSize
	subs, total, err := s.subRepo.ListSubs(ctx, pageSize, offset)
	if err != nil {
		s.logger.Errorf("ListSubs: %v", err)
		return nil, 0, err
	}
	s.logger.Infow("subscriptions listed", "page", page, "size", pageSize, "total", total)
	return subs, total, nil
}

func (s *subService) CalculateCost(ctx context.Context, filter models.PeriodFilter) (*models.CostResponse, error) {
	subs, err := s.subRepo.ListSubsForPeriod(ctx, filter)
	if err != nil {
		s.logger.Errorf("CalculateCost: %v", err)
		return nil, err
	}

	var details []models.CostDetail
	totalCost := 0

	for _, sub := range subs {
		overlapStart := sub.StartDate
		if filter.PeriodStart.After(overlapStart) {
			overlapStart = filter.PeriodStart
		}

		overlapEnd := filter.PeriodEnd
		if sub.EndDate != nil && sub.EndDate.Before(overlapEnd) {
			overlapEnd = *sub.EndDate
		}

		months := countMonths(overlapStart, overlapEnd)
		if months < 1 {
			months = 1
		}

		cost := int(sub.Price) * months
		totalCost += cost

		details = append(details, models.CostDetail{
			ID:           sub.ID,
			Name:         sub.Name,
			Price:        sub.Price,
			ActiveMonths: months,
			Cost:         cost,
		})
	}

	return &models.CostResponse{
		TotalCost: totalCost,
		Details:   details,
	}, nil
}

func countMonths(start, end time.Time) int {
	if start.After(end) {
		return 0
	}

	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	totalMonths := years*12 + months

	if end.Day() < start.Day() {
		totalMonths--
	}

	if totalMonths < 0 {
		totalMonths = 0
	}

	adjustedEnd := start.AddDate(0, totalMonths, 0)
	if end.After(adjustedEnd) || end.Equal(adjustedEnd) {
		totalMonths++
	}

	return totalMonths
}

func (s *subService) DeleteSubByID(ctx context.Context, id uuid.UUID) error {
	err := s.subRepo.DeleteSubByID(ctx, id)
	if err != nil {
		s.logger.Errorf("DeleteSubByID: %v", err)
		return err
	}
	s.logger.Infow("subscription deleted", "id", id)
	return nil
}
