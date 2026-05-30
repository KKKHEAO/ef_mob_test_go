package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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
	GetSubByID(w http.ResponseWriter, r *http.Request)
	UpdateSubByID(w http.ResponseWriter, r *http.Request)
	DeleteSubByID(w http.ResponseWriter, r *http.Request)
	ListSubs(w http.ResponseWriter, r *http.Request)
	CalculateCost(w http.ResponseWriter, r *http.Request)
}

// GetSubByID — GET /subscriptions/{id}
// @Summary      Получить подписку по ID
// @Description  Возвращает подписку по её идентификатору
// @Tags         subscriptions CRUDL
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      200  {object}  models.SubscriptionResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *subHandler) GetSubByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid uuid format"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	sub, err := h.subService.GetSubByID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{Error: "subscription not found"})
		return
	}

	resp := toResponse(sub)
	writeJSON(w, http.StatusOK, resp)
}

// CreateSub создаёт новую подписку
// @Summary      Создать подписку
// @Description  Создаёт новую запись о подписке пользователя
// @Tags         subscriptions CRUDL
// @Accept       json
// @Produce      json
// @Param        request body models.CreateSubscriptionRequest true "Данные подписки"
// @Success      201  {object}  models.SubscriptionResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions [post]
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

	newUuid, _ := uuid.NewV7()

	sub := &models.Subscription{
		ID:        newUuid,
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

// UpdateSubByID — PUT /subscriptions/{id}
// @Summary      Обновить подписку
// @Description  Обновляет название, стоимость и дату окончания подписки
// @Tags         subscriptions CRUDL
// @Accept       json
// @Produce      json
// @Param        id      path  string                        true  "UUID подписки"
// @Param        request body  models.UpdateSubscriptionRequest true  "Новые данные"
// @Success      200     {object}  models.SubscriptionResponse
// @Failure      400     {object}  models.ErrorResponse
// @Failure      404     {object}  models.ErrorResponse
// @Failure      500     {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *subHandler) UpdateSubByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid uuid format"})
		return
	}

	var req models.UpdateSubscriptionRequest
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

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "end_date must be in YYYY-MM-DD format"})
			return
		}
		endDate = &parsed
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	updated, err := h.subService.UpdateSubByID(ctx, &models.Subscription{
		ID:      id,
		Name:    req.Name,
		Price:   req.Price,
		EndDate: endDate,
	})
	if err != nil {
		writeJSON(w, http.StatusNotFound, models.ErrorResponse{Error: "subscription not found"})
		return
	}

	resp := toResponse(updated)
	writeJSON(w, http.StatusOK, resp)
}

// ListSubs — GET /subscriptions
// @Summary      Список подписок
// @Description  Возвращает список подписок с пагинацией
// @Tags         subscriptions CRUDL
// @Produce      json
// @Param        page      query  int  false  "Номер страницы (по умолчанию 1)"  default(1)
// @Param        page_size query  int  false  "Размер страницы (по умолчанию 10)"  default(10)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions [get]
func (h *subHandler) ListSubs(w http.ResponseWriter, r *http.Request) {
	// По дефолту
	page := 1
	pageSize := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	subs, total, err := h.subService.ListSubs(ctx, page, pageSize)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to list subscriptions"})
		return
	}

	data := make([]models.SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		data = append(data, toResponse(&sub))
	}

	totalPages := (total + pageSize - 1) / pageSize

	resp := models.PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	writeJSON(w, http.StatusOK, resp)
}

// CalculateCost — GET /subscriptions/cost
// @Summary      Стоимость подписок за период
// @Description  Считает суммарную стоимость подписок за указанный период, если у подписки нет окончания, то считается количество месяцев от старта, до конца даты расчета
// @Tags         subscriptions COST
// @Produce      json
// @Param        period_start  query  string  true   "Начало (YYYY-MM-DD)"
// @Param        period_end    query  string  true   "Конец (YYYY-MM-DD)"
// @Param        user_id       query  string  false  "Фильтр по ID пользователя"
// @Param        name          query  string  false  "Фильтр по названию"
// @Success      200  {object}  models.CostResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions/cost [get]
func (h *subHandler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	periodStartStr := r.URL.Query().Get("period_start")
	periodEndStr := r.URL.Query().Get("period_end")

	if periodStartStr == "" || periodEndStr == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "period_start and period_end are required"})
		return
	}

	periodStart, err := time.Parse("2006-01-02", periodStartStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "period_start must be YYYY-MM-DD"})
		return
	}

	periodEnd, err := time.Parse("2006-01-02", periodEndStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "period_end must be YYYY-MM-DD"})
		return
	}

	if !periodEnd.After(periodStart) {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "period_end must be after period_start"})
		return
	}

	filter := models.PeriodFilter{
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
	}

	if uidStr := r.URL.Query().Get("user_id"); uidStr != "" {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid user_id"})
			return
		}
		filter.UserID = &uid
	}

	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name = &name
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := h.subService.CalculateCost(ctx, filter)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to calculate cost"})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// DeleteSubByID — DELETE /subscriptions/{id}
// @Summary      Удалить подписку
// @Description  Удаляет подписку по её идентификатору
// @Tags         subscriptions CRUDL
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      204
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *subHandler) DeleteSubByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.ErrorResponse{Error: "invalid uuid format"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.subService.DeleteSubByID(ctx, id); err != nil {
		writeJSON(w, http.StatusInternalServerError, models.ErrorResponse{Error: "failed to delete subscription"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
