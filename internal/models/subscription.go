package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Price     int32      `json:"price" db:"price"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" db:"end_date"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionRequest struct {
	Name      string    `json:"name"`
	Price     int32     `json:"price"`
	UserID    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   *string   `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int32     `json:"price"`
	UserID    uuid.UUID `json:"user_id"`
	StartDate string    `json:"start_date"`
	EndDate   *string   `json:"end_date,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type UpdateSubscriptionRequest struct {
	Name    string  `json:"name"`
	Price   int32   `json:"price"`
	EndDate *string `json:"end_date,omitempty"`
}

type PaginationParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PaginatedResponse struct {
	Data       []SubscriptionResponse `json:"data"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

type PeriodFilter struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	PeriodStart time.Time  `json:"period_start"`
	PeriodEnd   time.Time  `json:"period_end"`
}

type CostDetail struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Price        int32     `json:"price"`
	ActiveMonths int       `json:"active_months"`
	Cost         int       `json:"cost"`
}

type CostResponse struct {
	TotalCost int          `json:"total_cost"`
	Details   []CostDetail `json:"details"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
