package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_api/modules/i18n/config"
	"github.com/khiemnd777/noah_api/modules/i18n/model"
	"github.com/khiemnd777/noah_api/modules/i18n/service"
	"github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/app/client_error"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/middleware/rbac"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils/table"
)

type LanguageHandler struct {
	svc  service.LanguageService
	deps *module.ModuleDeps[config.ModuleConfig]
}

func NewLanguageHandler(svc service.LanguageService, deps *module.ModuleDeps[config.ModuleConfig]) *LanguageHandler {
	return &LanguageHandler{svc: svc, deps: deps}
}

func (h *LanguageHandler) RegisterRoutes(router fiber.Router) {
	group := router.Group("/languages")
	app.RouterGet(group, "/", h.List)
	app.RouterGet(group, "/:language_id<int>", h.GetByID)
	app.RouterPost(group, "/", h.Create)
	app.RouterPut(group, "/:language_id<int>", h.Update)
	app.RouterDelete(group, "/:language_id<int>", h.Delete)
	app.RouterPost(group, "/:language_id<int>/import-xml", h.ImportXML)
	app.RouterGet(group, "/:language_id<int>/export-xml", h.ExportXML)
}

func (h *LanguageHandler) List(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.view"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	res, err := h.svc.List(c.UserContext(), table.ParseTableQuery(c, table.DefaultLimit))
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *LanguageHandler) GetByID(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.view"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	id, err := h.languageID(c)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "invalid id")
	}
	res, err := h.svc.GetByID(c.UserContext(), id)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *LanguageHandler) Create(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.create"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	var in model.LanguageDTO
	if err := c.BodyParser(&in); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	res, err := h.svc.Create(c.UserContext(), in)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *LanguageHandler) Update(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.update"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	id, err := h.languageID(c)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "invalid id")
	}
	var in model.LanguageDTO
	if err := c.BodyParser(&in); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	res, err := h.svc.Update(c.UserContext(), id, in)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *LanguageHandler) Delete(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.delete"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	id, err := h.languageID(c)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "invalid id")
	}
	if err := h.svc.Delete(c.UserContext(), id); err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

func (h *LanguageHandler) ImportXML(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.import"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	id, err := h.languageID(c)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "invalid id")
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "file is required")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "cannot open file")
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "cannot read file")
	}

	res, err := h.svc.ImportXML(c.UserContext(), id, body)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (h *LanguageHandler) ExportXML(c *fiber.Ctx) error {
	if err := h.guard(c, "languages.export"); err != nil {
		return client_error.ResponseError(c, fiber.StatusForbidden, err, err.Error())
	}
	id, err := h.languageID(c)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "invalid id")
	}
	doc, err := h.svc.ExportXML(c.UserContext(), id)
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	payload, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	filename := fmt.Sprintf("%s.xml", strings.ToLower(doc.Code))
	c.Set("Content-Type", "application/xml; charset=utf-8")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Send(append([]byte(xml.Header), payload...))
}

func (h *LanguageHandler) guard(c *fiber.Ctx, permissions ...string) error {
	return rbac.GuardAnyPermission(c, h.deps.Ent.(*generated.Client), permissions...)
}

func (h *LanguageHandler) languageID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("language_id"))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}
