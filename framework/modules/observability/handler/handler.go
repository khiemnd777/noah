package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/khiemnd777/noah_framework/modules/observability/model"
	"github.com/khiemnd777/noah_framework/modules/observability/service"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type PermissionGuard func(c frameworkhttp.Context, permission string) error

type Handler struct {
	svc               service.Service
	requirePermission PermissionGuard
}

func New(svc service.Service, requirePermission PermissionGuard) *Handler {
	return &Handler{svc: svc, requirePermission: requirePermission}
}

func (h *Handler) RegisterRoutes(router frameworkhttp.Router) {
	router.Get("/logs", h.ListLogs)
	router.Get("/logs/summary", h.GetSummary)
}

func (h *Handler) ListLogs(c frameworkhttp.Context) error {
	if err := h.guardPermission(c, "system_log.read"); err != nil {
		return err
	}

	query, err := parseListLogsQuery(c)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	items, err := h.svc.ListLogs(c.UserContext(), query)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err, "Failed to fetch system logs")
	}

	return c.JSON(map[string]any{
		"items": items,
		"total": len(items),
	})
}

func (h *Handler) GetSummary(c frameworkhttp.Context) error {
	if err := h.guardPermission(c, "system_log.read"); err != nil {
		return err
	}

	query, err := parseListLogsQuery(c)
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	summary, err := h.svc.GetSummary(c.UserContext(), query)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err, "Failed to fetch system log summary")
	}

	return c.JSON(summary)
}

func (h *Handler) guardPermission(c frameworkhttp.Context, permission string) error {
	if h.requirePermission == nil {
		return nil
	}
	return h.requirePermission(c, permission)
}

func parseListLogsQuery(c frameworkhttp.Context) (model.ListLogsQuery, error) {
	from, err := parseTimeQuery(c.Query("from"))
	if err != nil {
		return model.ListLogsQuery{}, err
	}
	to, err := parseTimeQuery(c.Query("to"))
	if err != nil {
		return model.ListLogsQuery{}, err
	}

	levels := splitCSV(c.Query("level"))
	if len(levels) == 0 {
		levels = splitCSV(c.Query("levels"))
	}

	userID, err := parseOptionalIntQuery(c, "user_id")
	if err != nil {
		return model.ListLogsQuery{}, err
	}
	deptID, err := parseOptionalIntQuery(c, "department_id")
	if err != nil {
		return model.ListLogsQuery{}, err
	}

	return model.ListLogsQuery{
		Levels:    levels,
		Module:    c.Query("module"),
		Service:   c.Query("service"),
		Env:       c.Query("env"),
		RequestID: c.Query("request_id"),
		UserID:    userID,
		DeptID:    deptID,
		Keyword:   c.Query("keyword"),
		Direction: c.Query("direction"),
		From:      from,
		To:        to,
		Limit:     c.QueryInt("limit", 50),
	}, nil
}

func parseTimeQuery(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, nil
	}

	if unixNanos, err := time.Parse(time.RFC3339Nano, raw); err == nil {
		return unixNanos, nil
	}

	if unixSeconds, err := time.Parse(time.RFC3339, raw); err == nil {
		return unixSeconds, nil
	}

	return time.Time{}, errors.New("invalid time format, expected RFC3339")
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func parseOptionalIntQuery(c frameworkhttp.Context, key string) (*int, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return nil, errors.New("invalid " + key)
	}
	return &value, nil
}

func respondError(c frameworkhttp.Context, statusCode int, err error, extraMessage string) error {
	message := "Server error"
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		message = "Client error"
	}
	if extraMessage != "" {
		message += ": " + extraMessage
	}
	if err != nil && extraMessage == "" {
		message += ": " + err.Error()
	}

	return c.Status(statusCode).JSON(map[string]any{
		"statusCode":    statusCode,
		"statusMessage": message,
	})
}
