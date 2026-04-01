package handler

import (
	"errors"
	"net/http"

	frameworkmodel "github.com/khiemnd777/noah_framework/modules/profile/model"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/profile/repository"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/profile/service"
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
	router.Put("/me/change-password", h.ChangePassword)
	router.Post("/me/exists-email", h.ExistsEmail)
	router.Post("/me/exists-phone", h.ExistsPhone)
	router.Get("/me", h.GetProfile)
	router.Put("/me", h.UpdateProfile)
	router.Delete("/me", h.Delete)
}

func (h *Handler) GetProfile(c frameworkhttp.Context) error {
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	profile, err := h.svc.GetProfile(c.UserContext(), userID)
	if err != nil {
		if err == frameworkrepository.ErrNotFound {
			return respondError(c, http.StatusNotFound, "User not found")
		}
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(profile)
}

func (h *Handler) ExistsEmail(c frameworkhttp.Context) error {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	exists, err := h.svc.CheckEmailExists(c.UserContext(), userID, body.Email)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(exists)
}

func (h *Handler) ExistsPhone(c frameworkhttp.Context) error {
	var body struct {
		Phone string `json:"phone"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	exists, err := h.svc.CheckPhoneExists(c.UserContext(), userID, body.Phone)
	if err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.JSON(exists)
}

func (h *Handler) UpdateProfile(c frameworkhttp.Context) error {
	var body struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Avatar string  `json:"avatar"`
		Phone  *string `json:"phone"`
		Email  *string `json:"email"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	updated, err := h.svc.UpdateProfile(c.UserContext(), userID, body.Name, body.Avatar, body.Phone, body.Email)
	if err != nil {
		switch {
		case errors.Is(err, frameworkmodel.ErrEmailExists), errors.Is(err, frameworkmodel.ErrPhoneExists):
			return c.Status(http.StatusOK).JSON(map[string]any{
				"statusCode":    102,
				"errorCode":     "ErrEmailOrPhoneExists",
				"statusMessage": err.Error(),
			})
		case errors.Is(err, frameworkrepository.ErrNotFound):
			return respondError(c, http.StatusNotFound, "User not found")
		default:
			return respondError(c, http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(updated)
}

func (h *Handler) ChangePassword(c frameworkhttp.Context) error {
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondError(c, http.StatusBadRequest, "Invalid request body")
	}
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	if err := h.svc.ChangePassword(c.UserContext(), userID, body.CurrentPassword, body.NewPassword); err != nil {
		return respondError(c, http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusAccepted)
}

func (h *Handler) Delete(c frameworkhttp.Context) error {
	userID, ok := frameworkruntime.UserIDFromContext(c)
	if !ok {
		return respondError(c, http.StatusUnauthorized, "Unauthorized")
	}
	if err := h.svc.Delete(c.UserContext(), userID); err != nil {
		if err == frameworkrepository.ErrNotFound {
			return respondError(c, http.StatusNotFound, "User not found")
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
