package handler

import (
	"net/http"
	"strconv"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/folder/model"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/folder/repository"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/folder/service"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type Handler struct {
	svc *frameworkservice.Service
}

func New(svc *frameworkservice.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(router frameworkhttp.Router) {
	router.Post("/", h.Create)
	router.Get("/", h.List)
	router.Get("/:id", h.Get)
	router.Put("/:id", h.Update)
	router.Delete("/:id", h.Delete)
}

func (h *Handler) Create(c frameworkhttp.Context) error {
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	var input frameworkmodel.Folder
	if err := c.BodyParser(&input); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid input")
	}
	folder, err := h.svc.Create(c.UserContext(), userID, input)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(folder)
}

func (h *Handler) Get(c frameworkhttp.Context) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid folder id")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	folder, err := h.svc.Get(c.UserContext(), id, userID)
	if err != nil {
		if err == frameworkrepository.ErrNotFound {
			return respondError(c, http.StatusNotFound, "Folder not found")
		}
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(folder)
}

func (h *Handler) List(c frameworkhttp.Context) error {
	userID := c.QueryInt("user_id", 0)
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	items, hasMore, err := h.svc.ListPaginated(c.UserContext(), userID, page, limit)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(map[string]any{
		"items":   items,
		"hasMore": hasMore,
	})
}

func (h *Handler) Update(c frameworkhttp.Context) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid folder id")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	var input frameworkmodel.Folder
	if err := c.BodyParser(&input); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid input")
	}
	updated, err := h.svc.Update(c.UserContext(), id, userID, input)
	if err != nil {
		if err == frameworkrepository.ErrNotFound {
			return respondError(c, http.StatusNotFound, "Folder not found")
		}
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(updated)
}

func (h *Handler) Delete(c frameworkhttp.Context) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid folder id")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	if err := h.svc.Delete(c.UserContext(), id, userID); err != nil {
		if err == frameworkrepository.ErrNotFound {
			return respondError(c, http.StatusNotFound, "Folder not found")
		}
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusNoContent)
}

func respondError(c frameworkhttp.Context, statusCode int, message string) error {
	return c.Status(statusCode).JSON(map[string]any{
		"statusCode":    statusCode,
		"statusMessage": message,
	})
}
