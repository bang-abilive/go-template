package handlers

import (
	"net/http"

	"ndinhbang/go-template/pkg/authorizer"
	"ndinhbang/go-template/pkg/db"
	"ndinhbang/go-template/pkg/server"

	"github.com/labstack/echo/v5"
)

var _ server.GroupRegistrar = (*AuthorizeHandler)(nil)

// AuthorizeHandler exposes a debug endpoint to check ABAC policy decisions.
// It does NOT sit behind the auth/casbin middlewares, so any user_id (or none) can be tested.
type AuthorizeHandler struct {
	auth *authorizer.Authorizer
	db   *db.PostgresDatabase
}

func NewAuthorizeHandler(auth *authorizer.Authorizer, database *db.PostgresDatabase) *AuthorizeHandler {
	return &AuthorizeHandler{auth: auth, db: database}
}

func (h *AuthorizeHandler) RegisterRoutes(g *echo.Group) {
	g.AddRoute(echo.Route{
		Method:  "POST",
		Path:    "/authorize",
		Handler: h.Check,
		Name:    "authorize.check",
	})
}

type authorizeRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	Object string `json:"object"  validate:"required"`
	Action string `json:"action"  validate:"required"`
}

// Check evaluates an ABAC policy for the supplied (user_id, object, action) triple.
func (h *AuthorizeHandler) Check(c *echo.Context) error {
	var req authorizeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	// Resolve UserAttr for the given user_id.
	var attr authorizer.UserAttr
	attr.ID = req.UserID
	err := h.db.Pool().QueryRow(
		c.Request().Context(),
		`SELECT r.slug, r.lv
		   FROM roles r
		   JOIN user_role ur ON ur.role_id = r.id
		  WHERE ur.user_id = $1
		  ORDER BY r.lv DESC
		  LIMIT 1`,
		req.UserID,
	).Scan(&attr.Role, &attr.Level)
	if err != nil {
		// User has no role: treat as anonymous with zero attrs.
		attr.Role = ""
		attr.Level = 0
	}

	allowed, err := h.auth.GetEnforcer().Enforce(attr, req.Object, req.Action)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "enforce error: " + err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"allowed": allowed,
		"user": map[string]any{
			"id":    attr.ID,
			"role":  attr.Role,
			"level": attr.Level,
		},
		"object": req.Object,
		"action": req.Action,
	})
}
