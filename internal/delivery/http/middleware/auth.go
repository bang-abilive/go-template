package middleware

import (
	"net/http"
	"strconv"

	"ndinhbang/go-template/pkg/authorizer"
	"ndinhbang/go-template/pkg/db"

	"github.com/labstack/echo/v5"
)

// AuthMiddleware resolves a UserAttr from ?user_id query param and stores it on the request context.
// This is a mock auth layer — no JWT or session verification is performed.
type AuthMiddleware struct {
	db *db.PostgresDatabase
}

func NewAuthMiddleware(database *db.PostgresDatabase) *AuthMiddleware {
	return &AuthMiddleware{db: database}
}

func (m *AuthMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			rawID := c.QueryParam("user_id")
			if rawID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing user_id"})
			}
			userID, err := strconv.ParseInt(rawID, 10, 64)
			if err != nil || userID <= 0 {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
			}

			// Fetch the highest-level role for this user.
			var attr authorizer.UserAttr
			attr.ID = userID
			err = m.db.Pool().QueryRow(
				c.Request().Context(),
				`SELECT r.slug, r.lv
				   FROM roles r
				   JOIN user_role ur ON ur.role_id = r.id
				  WHERE ur.user_id = $1
				  ORDER BY r.lv DESC
				  LIMIT 1`,
				userID,
			).Scan(&attr.Role, &attr.Level)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found or has no role"})
			}

			c.SetRequest(c.Request().WithContext(
				authorizer.WithContext(c.Request().Context(), attr),
			))
			return next(c)
		}
	}
}
