package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_framework/modules/token/service"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app/client_error"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type TokenHandler struct {
	svc *service.TokenService
}

func NewTokenHandler(svc *service.TokenService) *TokenHandler {
	return &TokenHandler{svc: svc}
}

func (h *TokenHandler) RegisterRoutes(router frameworkhttp.Router) {
	app.RouterGet(router, "/", h.Token)
	app.RouterPost(router, "/refresh", h.RefreshTokens)
	app.RouterPost(router, "/generate", h.GenerateTokens)
	app.RouterDelete(router, "/delete", h.DeleteRefreshToken)
}

func (h *TokenHandler) Token(c frameworkhttp.Context) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TokenHandler) RefreshTokens(c frameworkhttp.Context) error {
	type req struct {
		Token string `json:"refreshToken"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "Invalid request")
	}
	tokens, err := h.svc.RefreshToken(c.UserContext(), body.Token)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			return client_error.ResponseError(c, fiber.StatusUnauthorized, err, err.Error())
		}
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}

	return c.JSON(tokens)
}

func (h *TokenHandler) GenerateTokens(c frameworkhttp.Context) error {
	type req struct {
		UserID int `json:"userID"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "Invalid request")
	}

	tokens, err := h.svc.GenerateTokens(c.UserContext(), body.UserID)

	if err != nil {
		return client_error.ResponseError(c, fiber.StatusUnauthorized, err, err.Error())
	}

	return c.JSON(tokens)
}

func (h *TokenHandler) DeleteRefreshToken(c frameworkhttp.Context) error {
	type req struct {
		RefreshToken string `json:"refreshToken"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return client_error.ResponseError(c, fiber.StatusBadRequest, err, "Invalid request")
	}
	if err := h.svc.DeleteRefreshToken(c.UserContext(), body.RefreshToken); err != nil {
		return client_error.ResponseError(c, fiber.StatusInternalServerError, err, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
