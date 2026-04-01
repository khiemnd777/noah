// scripts/create_module/templates/handler_http.go.tmpl
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_api/modules/attribute/config"
	"github.com/khiemnd777/noah_api/modules/attribute/service"
	"github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/app/client_error"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type AttributeHandler struct {
	svc  *service.AttributeService
	deps *module.ModuleDeps[config.ModuleConfig]
}

func NewAttributeHandler(svc *service.AttributeService, deps *module.ModuleDeps[config.ModuleConfig]) *AttributeHandler {
	return &AttributeHandler{
		svc:  svc,
		deps: deps,
	}
}

func (h *AttributeHandler) RegisterRoutes(router frameworkhttp.Router) {
	app.RouterGet(router, "/", h.List)
	app.RouterPost(router, "/", h.Create)
	app.RouterGet(router, "/:id", h.Get)
	app.RouterPut(router, "/:id", h.Update)
	app.RouterDelete(router, "/:id", h.Delete)
}

func (h *AttributeHandler) Create(c frameworkhttp.Context) error {
	var req struct {
		UserID        int    `json:"user_id"`
		AttributeName string `json:"attribute_name"`
		AttributeType string `json:"attribute_type"`
	}
	if err := c.BodyParser(&req); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	result, err := h.svc.CreateAttribute(c.UserContext(), req.UserID, req.AttributeName, req.AttributeType)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeHandler) Get(c frameworkhttp.Context) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userId, _ := utils.GetUserIDInt(c)
	result, err := h.svc.GetAttribute(c.UserContext(), id, userId)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusNotFound, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeHandler) List(c frameworkhttp.Context) error {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	items, hasMore, err := h.svc.ListAttributesPaginated(c.UserContext(), userId, page, limit)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(fiber.Map{
		"items":   items,
		"hasMore": hasMore,
	})
}

func (h *AttributeHandler) Update(c frameworkhttp.Context) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userId, _ := utils.GetUserIDInt(c)
	var req struct {
		AttributeName string `json:"attribute_name"`
		AttributeType string `json:"attribute_type"`
	}
	if err := c.BodyParser(&req); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	result, err := h.svc.UpdateAttribute(c.UserContext(), id, userId, req.AttributeName, req.AttributeType)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.JSON(result)
}

func (h *AttributeHandler) Delete(c frameworkhttp.Context) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userId, _ := utils.GetUserIDInt(c)
	if err := h.svc.DeleteAttribute(c.UserContext(), id, userId); err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
