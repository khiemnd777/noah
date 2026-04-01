package handler

import (
	"errors"
	"net/http"
	"strconv"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/attribute/model"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/attribute/repository"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/attribute/service"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type AttributeHandler struct {
	svc *frameworkservice.AttributeService
}

func NewAttributeHandler(svc *frameworkservice.AttributeService) *AttributeHandler {
	return &AttributeHandler{svc: svc}
}

func (h *AttributeHandler) RegisterRoutes(router frameworkhttp.Router) {
	router.Get("/", h.List)
	router.Post("/", h.Create)
	router.Get("/:id", h.Get)
	router.Put("/:id", h.Update)
	router.Delete("/:id", h.Delete)
}

func (h *AttributeHandler) Create(c frameworkhttp.Context) error {
	var req struct {
		UserID        int    `json:"user_id"`
		AttributeName string `json:"attribute_name"`
		AttributeType string `json:"attribute_type"`
	}
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	result, err := h.svc.CreateAttribute(c.UserContext(), req.UserID, req.AttributeName, req.AttributeType)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeHandler) Get(c frameworkhttp.Context) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	result, err := h.svc.GetAttribute(c.UserContext(), id, userID)
	if err != nil {
		if frameworkrepository.IsNotFound(err) {
			return respondError(c, http.StatusNotFound, err, err.Error())
		}
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeHandler) List(c frameworkhttp.Context) error {
	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		userID = 0
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	items, hasMore, err := h.svc.ListAttributesPaginated(c.UserContext(), userID, page, limit)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(map[string]any{
		"items":   items,
		"hasMore": hasMore,
	})
}

func (h *AttributeHandler) Update(c frameworkhttp.Context) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	var req struct {
		AttributeName string `json:"attribute_name"`
		AttributeType string `json:"attribute_type"`
	}
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	result, err := h.svc.UpdateAttribute(c.UserContext(), id, userID, req.AttributeName, req.AttributeType)
	if err != nil {
		if frameworkrepository.IsNotFound(err) {
			return respondError(c, http.StatusNotFound, err, err.Error())
		}
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeHandler) Delete(c frameworkhttp.Context) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	if err := h.svc.DeleteAttribute(c.UserContext(), id, userID); err != nil {
		if frameworkrepository.IsNotFound(err) {
			return respondError(c, http.StatusNotFound, err, err.Error())
		}
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

type AttributeOptionHandler struct {
	svc *frameworkservice.AttributeOptionService
}

func NewAttributeOptionHandler(svc *frameworkservice.AttributeOptionService) *AttributeOptionHandler {
	return &AttributeOptionHandler{svc: svc}
}

func (h *AttributeOptionHandler) RegisterRoutes(router frameworkhttp.Router) {
	router.Route("/:id/options", func(r frameworkhttp.Router) {
		r.Get("/", h.ListOptions)
		r.Post("/", h.CreateOption)
		r.Put("/reorder", h.ReorderOptions)
		r.Put("/:option_id", h.UpdateOption)
		r.Delete("/:option_id", h.DeleteOption)
	})
}

func (h *AttributeOptionHandler) ListOptions(c frameworkhttp.Context) error {
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	result, err := h.svc.ListOptions(c.UserContext(), userID, attributeID)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeOptionHandler) CreateOption(c frameworkhttp.Context) error {
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	var req struct {
		OptionValue  string `json:"option_value"`
		DisplayOrder int    `json:"display_order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	result, err := h.svc.CreateOption(c.UserContext(), userID, attributeID, req.OptionValue, req.DisplayOrder)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeOptionHandler) UpdateOption(c frameworkhttp.Context) error {
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	optionID, err := strconv.Atoi(c.Params("option_id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid option id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	var req struct {
		OptionValue  string `json:"option_value"`
		DisplayOrder int    `json:"display_order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	result, err := h.svc.UpdateOption(c.UserContext(), userID, attributeID, optionID, req.OptionValue, req.DisplayOrder)
	if err != nil {
		if frameworkrepository.IsNotFound(err) {
			return respondError(c, http.StatusNotFound, err, err.Error())
		}
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeOptionHandler) DeleteOption(c frameworkhttp.Context) error {
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	optionID, err := strconv.Atoi(c.Params("option_id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid option id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	if err := h.svc.DeleteOption(c.UserContext(), userID, attributeID, optionID); err != nil {
		if frameworkrepository.IsNotFound(err) {
			return respondError(c, http.StatusNotFound, err, err.Error())
		}
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

func (h *AttributeOptionHandler) ReorderOptions(c frameworkhttp.Context) error {
	attributeID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, err, "invalid attribute id")
	}
	userID, err := userIDFromContext(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, err, err.Error())
	}

	var req struct {
		Orders []frameworkmodel.OptionOrder `json:"orders"`
	}
	if err := c.BodyParser(&req); err != nil {
		return respondError(c, http.StatusBadRequest, err, err.Error())
	}

	if err := h.svc.BatchUpdateDisplayOrder(c.UserContext(), userID, attributeID, req.Orders); err != nil {
		return respondError(c, http.StatusInternalServerError, err, err.Error())
	}
	return c.SendStatus(http.StatusOK)
}

func userIDFromContext(c frameworkhttp.Context) (int, error) {
	switch value := c.Locals("userID").(type) {
	case int:
		return value, nil
	case int64:
		return int(value), nil
	case float64:
		return int(value), nil
	case string:
		userID, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.New("invalid user context")
		}
		return userID, nil
	default:
		return 0, errors.New("missing user context")
	}
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
