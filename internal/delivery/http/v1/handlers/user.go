package handlers

import (
	"ndinhbang/go-template/internal/usecase/user"
	"ndinhbang/go-template/pkg/server"
	"net/http"

	"github.com/labstack/echo/v5"
)

var _ server.GroupRegistrar = (*UserHandler)(nil)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(svc user.Service) *UserHandler {
	return &UserHandler{service: svc}
}

func (h *UserHandler) RegisterRoutes(e *echo.Group) {
	g := e.Group("/user")
	g.AddRoute(echo.Route{Method: "POST", Path: "/register", Handler: h.Register, Name: "user.register"})
}

func (h *UserHandler) Register(ctx *echo.Context) error {
	var in user.RegisterRequest
	if err := ctx.Bind(&in); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	out, err := h.service.Register(ctx.Request().Context(), in)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusCreated, out)
}
