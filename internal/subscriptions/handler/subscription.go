package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"ef_mob_test_go/internal/models"
	"ef_mob_test_go/internal/subscriptions/service"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type subHandler struct {
	subService service.SubService
	logger     *zap.SugaredLogger
}

func NewSubHandler(subService service.SubService, logger *zap.SugaredLogger) SubHandler {
	return &subHandler{
		subService: subService,
		logger:     logger,
	}
}

type SubHandler interface {
	CreateSub(w http.ResponseWriter, r *http.Request)
}

func (h *subHandler) CreateSub(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "name is required"})
		return
	}
	if req.Price < 0 {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "price must be >= 0"})
		return
	}
	if req.UserID == uuid.Nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "user_id is required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "start_date must be in YYYY-MM-DD format"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "end_date must be in YYYY-MM-DD format"})
			return
		}
		endDate = &parsed

		if !endDate.After(startDate) {
			writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "end_date must be after start_date"})
			return
		}
	}

	sub := &models.Subscription{
		ID:        uuid.New(),
		Name:      req.Name,
		Price:     req.Price,
		UserID:    req.UserID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	created, err := h.subService.CreateSub(ctx, sub)
	if err != nil {
		h.logger.Errorf("CreateSub failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to create subscription"})
		return
	}

	resp := toResponse(created)
	writeJSON(w, http.StatusCreated, resp)
}

// TODO: вынести в helpers.go
func toResponse(sub *models.Subscription) models.SubscriptionResponse {
	var endDate *string
	if sub.EndDate != nil {
		s := sub.EndDate.Format("2006-01-02")
		endDate = &s
	}

	return models.SubscriptionResponse{
		ID:        sub.ID,
		Name:      sub.Name,
		Price:     sub.Price,
		UserID:    sub.UserID,
		StartDate: sub.StartDate.Format("2006-01-02"),
		EndDate:   endDate,
		CreatedAt: sub.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sub.UpdatedAt.Format(time.RFC3339),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg})
}
