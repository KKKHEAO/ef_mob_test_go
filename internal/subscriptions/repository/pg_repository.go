package repository

import (
	"context"
	"database/sql"
	"ef_mob_test_go/internal/models"
	"errors"
)

type subRepository struct {
	db *sql.DB
}

type SubRepo interface {
	CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
}

func NewSubRepository(db *sql.DB) SubRepo {
	return &subRepository{db: db}
}

func (r *subRepository) CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	u := &models.Subscription{}

	err := r.db.QueryRowContext(
		ctx, createSubscription,
		sub.ID, sub.Name, sub.Price, sub.UserID, sub.StartDate, sub.EndDate,
	).Scan(
		&u.ID, &u.Name, &u.Price, &u.UserID,
		&u.StartDate, &u.EndDate, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("Error during create subscription")
	}

	return u, nil
}
