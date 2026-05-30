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

type ErrorResponse struct {
	Error string `json:"error"`
}
