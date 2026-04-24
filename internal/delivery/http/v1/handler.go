package v1

import (
	"ndinhbang/go-template/internal/usecase/user"
	"net/http"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(svc user.Service) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) Register(ctx *echo.Context) error {
	var in user.RegisterUserRequest
	if err := ctx.Bind(&in); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	out, err := h.service.Register(ctx.Request().Context(), in)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, out)
}
