package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	moduleModel "github.com/khiemnd777/noah_framework/modules/observability/model"
)

type LokiConfig struct {
	BaseURL        string
	TenantID       string
	BearerToken    string
	Timeout        time.Duration
	StreamSelector string
	MaxQueryLimit  int
}

type Repository interface {
	ListLogs(ctx context.Context, query moduleModel.ListLogsQuery) ([]moduleModel.LogEntry, error)
	GetSummary(ctx context.Context, query moduleModel.ListLogsQuery) (moduleModel.Summary, error)
}

type observabilityRepository struct {
	client *http.Client
	cfg    LokiConfig
}

func New(cfg LokiConfig) Repository {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &observabilityRepository{
		client: &http.Client{Timeout: timeout},
		cfg:    cfg,
	}
}

func (r *observabilityRepository) ListLogs(ctx context.Context, query moduleModel.ListLogsQuery) ([]moduleModel.LogEntry, error) {
	end := query.To
	if end.IsZero() {
		end = time.Now()
	}
	start := query.From
	if start.IsZero() {
		start = end.Add(-1 * time.Hour)
	}

	params := url.Values{}
	params.Set("query", r.buildStreamQuery(query))
	params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	params.Set("limit", strconv.Itoa(r.normalizeLimit(query.Limit)))
	params.Set("direction", normalizeDirection(query.Direction))

	var resp lokiRangeResponse
	if err := r.doRequest(ctx, "/loki/api/v1/query_range", params, &resp); err != nil {
		return nil, err
	}

	items := make([]moduleModel.LogEntry, 0, r.normalizeLimit(query.Limit))
	for _, stream := range resp.Data.Result {
		for _, value := range stream.Values {
			if len(value) < 2 {
				continue
			}
			entry := parseLogEntry(value[0], value[1], stream.Stream)
			items = append(items, entry)
		}
	}

	return items, nil
}

func (r *observabilityRepository) GetSummary(ctx context.Context, query moduleModel.ListLogsQuery) (moduleModel.Summary, error) {
	end := query.To
	if end.IsZero() {
		end = time.Now()
	}
	start := query.From
	if start.IsZero() {
		start = end.Add(-1 * time.Hour)
	}

	summary := moduleModel.Summary{
		From: start,
		To:   end,
		Counts: map[string]int{
			"warn":  0,
			"error": 0,
		},
	}

	for _, level := range []string{"warn", "error"} {
		value, err := r.queryCount(ctx, query, level, end, end.Sub(start))
		if err != nil {
			return moduleModel.Summary{}, err
		}
		summary.Counts[level] = value
	}

	return summary, nil
}

func (r *observabilityRepository) queryCount(ctx context.Context, query moduleModel.ListLogsQuery, level string, at time.Time, window time.Duration) (int, error) {
	if window <= 0 {
		window = time.Hour
	}

	params := url.Values{}
	params.Set("query", fmt.Sprintf(`sum(count_over_time((%s)[%s]))`, r.buildStreamQuery(withLevels(query, level)), formatRange(window)))
	params.Set("time", strconv.FormatInt(at.UnixNano(), 10))

	var resp lokiInstantResponse
	if err := r.doRequest(ctx, "/loki/api/v1/query", params, &resp); err != nil {
		return 0, err
	}

	if len(resp.Data.Result) == 0 || len(resp.Data.Result[0].Value) < 2 {
		return 0, nil
	}

	rawCount, _ := resp.Data.Result[0].Value[1].(string)
	count, err := strconv.ParseFloat(rawCount, 64)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *observabilityRepository) buildStreamQuery(query moduleModel.ListLogsQuery) string {
	selector := strings.TrimSpace(r.cfg.StreamSelector)
	if selector == "" {
		selector = `{app="noah_api"}`
	}

	parts := []string{selector}
	if query.Keyword != "" {
		parts = append(parts, fmt.Sprintf(`|= %s`, strconv.Quote(query.Keyword)))
	}

	parts = append(parts, "| json")
	parts = append(parts, `| __error__=""`)

	if len(query.Levels) > 0 {
		if len(query.Levels) == 1 {
			parts = append(parts, fmt.Sprintf(`| level=%s`, strconv.Quote(strings.ToLower(query.Levels[0]))))
		} else {
			levels := make([]string, 0, len(query.Levels))
			for _, level := range query.Levels {
				levels = append(levels, regexp.QuoteMeta(strings.ToLower(level)))
			}
			parts = append(parts, fmt.Sprintf(`| level=~%s`, strconv.Quote(strings.Join(levels, "|"))))
		}
	}
	if query.Module != "" {
		parts = append(parts, fmt.Sprintf(`| module=%s`, strconv.Quote(query.Module)))
	}
	if query.Service != "" {
		parts = append(parts, fmt.Sprintf(`| service=%s`, strconv.Quote(query.Service)))
	}
	if query.Env != "" {
		parts = append(parts, fmt.Sprintf(`| env=%s`, strconv.Quote(query.Env)))
	}
	if query.RequestID != "" {
		parts = append(parts, fmt.Sprintf(`| request_id=%s`, strconv.Quote(query.RequestID)))
	}
	if query.UserID != nil {
		parts = append(parts, fmt.Sprintf(`| user_id=%s`, strconv.Quote(strconv.Itoa(*query.UserID))))
	}
	if query.DeptID != nil {
		parts = append(parts, fmt.Sprintf(`| department_id=%s`, strconv.Quote(strconv.Itoa(*query.DeptID))))
	}

	return strings.Join(parts, " ")
}

func (r *observabilityRepository) doRequest(ctx context.Context, path string, params url.Values, out any) error {
	baseURL := strings.TrimRight(r.cfg.BaseURL, "/")
	if baseURL == "" {
		return fmt.Errorf("observability loki base_url is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}

	if r.cfg.TenantID != "" {
		req.Header.Set("X-Scope-OrgID", r.cfg.TenantID)
	}
	if r.cfg.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+r.cfg.BearerToken)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("loki request failed: %s", strings.TrimSpace(string(body)))
	}

	return json.Unmarshal(body, out)
}

func parseLogEntry(rawTS, rawLine string, stream map[string]string) moduleModel.LogEntry {
	entry := moduleModel.LogEntry{
		Raw: rawLine,
	}

	if ts, err := strconv.ParseInt(rawTS, 10, 64); err == nil && ts > 0 {
		entry.Timestamp = time.Unix(0, ts).UTC()
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(rawLine), &payload); err == nil {
		entry.Level = normalizeStringField(payload, "level")
		entry.Message = normalizeStringField(payload, "message")
		entry.Service = normalizeStringField(payload, "service")
		entry.Module = normalizeStringField(payload, "module")
		entry.Environment = normalizeStringField(payload, "env")
		entry.RequestID = normalizeStringField(payload, "request_id")
		entry.Method = normalizeStringField(payload, "method")
		entry.Path = normalizeStringField(payload, "path")
		entry.Source = normalizeStringField(payload, "source")
		entry.Error = normalizeStringField(payload, "error")
		entry.Stacktrace = normalizeStringField(payload, "stacktrace")

		if tsValue := normalizeStringField(payload, "ts"); tsValue != "" {
			if parsed, err := time.Parse(time.RFC3339Nano, tsValue); err == nil {
				entry.Timestamp = parsed.UTC()
			}
		}

		entry.UserID = normalizeIntField(payload, "user_id")
		entry.DepartmentID = normalizeIntField(payload, "department_id")

		delete(payload, "ts")
		delete(payload, "level")
		delete(payload, "message")
		delete(payload, "service")
		delete(payload, "module")
		delete(payload, "env")
		delete(payload, "request_id")
		delete(payload, "user_id")
		delete(payload, "department_id")
		delete(payload, "method")
		delete(payload, "path")
		delete(payload, "source")
		delete(payload, "error")
		delete(payload, "stacktrace")
		if len(payload) > 0 {
			entry.Fields = payload
		}
	}

	if entry.Service == "" {
		entry.Service = stream["service"]
	}
	if entry.Module == "" {
		entry.Module = stream["module"]
	}
	if entry.Environment == "" {
		entry.Environment = stream["env"]
	}
	if entry.RequestID == "" {
		entry.RequestID = stream["request_id"]
	}
	if entry.UserID == 0 {
		entry.UserID = parseStringInt(stream["user_id"])
	}
	if entry.DepartmentID == 0 {
		entry.DepartmentID = parseStringInt(stream["department_id"])
	}

	return entry
}

func normalizeStringField(payload map[string]any, key string) string {
	value, ok := payload[key]
	if !ok {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}

func normalizeIntField(payload map[string]any, key string) int {
	value, ok := payload[key]
	if !ok {
		return 0
	}
	switch typed := value.(type) {
	case int:
		return typed
	case int64:
		return int(typed)
	case float64:
		return int(typed)
	case string:
		return parseStringInt(typed)
	default:
		return 0
	}
}

func parseStringInt(raw string) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0
	}
	return value
}

func (r *observabilityRepository) normalizeLimit(limit int) int {
	if limit <= 0 {
		limit = 50
	}
	if r.cfg.MaxQueryLimit > 0 && limit > r.cfg.MaxQueryLimit {
		return r.cfg.MaxQueryLimit
	}
	return limit
}

func normalizeDirection(direction string) string {
	switch strings.ToUpper(strings.TrimSpace(direction)) {
	case "FORWARD":
		return "FORWARD"
	default:
		return "BACKWARD"
	}
}

func withLevels(query moduleModel.ListLogsQuery, level string) moduleModel.ListLogsQuery {
	query.Levels = []string{level}
	return query
}

func formatRange(window time.Duration) string {
	if window%time.Hour == 0 {
		return fmt.Sprintf("%dh", int(window/time.Hour))
	}
	if window%time.Minute == 0 {
		return fmt.Sprintf("%dm", int(window/time.Minute))
	}
	return window.String()
}

type lokiRangeResponse struct {
	Data struct {
		Result []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type lokiInstantResponse struct {
	Data struct {
		Result []struct {
			Value []any `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
