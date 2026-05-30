package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	ListSubsForPeriod(ctx context.Context, filter models.PeriodFilter) ([]models.Subscription, error)
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
		return nil, fmt.Errorf("create subscription: %w", err)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get subscription by id: %w", models.ErrNotFound)
		}
		return nil, fmt.Errorf("get subscription by id: %w", err)
	}

	return u, nil
}

func (r *subRepository) DeleteSubByID(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, deleteSubscriptionByID, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("delete subscription: %w", models.ErrNotFound)
	}

	return nil
}

func (r *subRepository) ListSubs(ctx context.Context, limit, offset int) ([]models.Subscription, int, error) {
	rows, err := r.db.QueryContext(ctx, listSubscriptions, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Price, &sub.UserID,
			&sub.StartDate, &sub.EndDate, &sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration: %w", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countSubscriptions).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count subscriptions: %w", err)
	}

	return subs, total, nil
}

func (r *subRepository) ListSubsForPeriod(ctx context.Context, filter models.PeriodFilter) ([]models.Subscription, error) {
	query := listSubsForPeriod
	args := make([]any, 0)

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", len(args)+3)
		args = append(args, *filter.UserID)
	}

	if filter.Name != nil {
		query += fmt.Sprintf(" AND name = $%d", len(args)+3)
		args = append(args, *filter.Name)
	}

	allArgs := []any{filter.PeriodEnd, filter.PeriodStart}
	allArgs = append(allArgs, args...)

	rows, err := r.db.QueryContext(ctx, query, allArgs...)
	if err != nil {
		return nil, fmt.Errorf("list subs for period: %w", err)
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(
			&sub.ID, &sub.Name, &sub.Price, &sub.UserID,
			&sub.StartDate, &sub.EndDate,
		); err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return subs, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("update subscription: %w", models.ErrNotFound)
		}
		return nil, fmt.Errorf("update subscription: %w", err)
	}

	return u, nil
}
