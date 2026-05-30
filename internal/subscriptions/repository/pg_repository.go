package repository

import (
	"context"
	"database/sql"

	"ef_mob_test_go/internal/models"
	"errors"

	"github.com/google/uuid"
)

type subRepository struct {
	db *sql.DB
}

type SubRepo interface {
	CreateSub(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	GetSubByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	DeleteSubByID(ctx context.Context, id uuid.UUID) error
	UpdateSubByID(ctx context.Context, sub *models.Subscription) (*models.Subscription, error)
	ListSubs(ctx context.Context, limit, offset int) ([]models.Subscription, int, error)
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

func (r *subRepository) GetSubByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	u := &models.Subscription{}

	err := r.db.QueryRowContext(ctx, getSubscriptionByID, id).Scan(
		&u.ID, &u.Name, &u.Price, &u.UserID,
		&u.StartDate, &u.EndDate, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("Error during get suscription by id")
	}

	return u, nil
}

func (r *subRepository) DeleteSubByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteSubscriptionByID, id)
	if err != nil {
		return errors.New("Error during delete suscription by id")
	}

	return nil
}

func (r *subRepository) ListSubs(ctx context.Context, limit, offset int) ([]models.Subscription, int, error) {
	rows, err := r.db.QueryContext(ctx, listSubscriptions, limit, offset)
	if err != nil {
		return nil, 0, errors.New("Error during list subscriptions")
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Price, &sub.UserID,
			&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, 0, errors.New("Error scanning subscription")
		}
		subs = append(subs, sub)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSubscriptions).Scan(&total); err != nil {
		return nil, 0, errors.New("Error counting subscriptions")
	}

	return subs, total, nil
}

func (r *subRepository) UpdateSubByID(ctx context.Context, sub *models.Subscription) (*models.Subscription, error) {
	u := &models.Subscription{}

	err := r.db.QueryRowContext(
		ctx, updateSubscriptionByID,
		sub.ID, sub.Name, sub.Price, sub.EndDate,
	).Scan(
		&u.ID, &u.Name, &u.Price, &u.UserID,
		&u.StartDate, &u.EndDate, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, errors.New("Error during update subscription")
	}

	return u, nil
}
