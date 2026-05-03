package middleware

import (
	"net/http"

	"ndinhbang/go-template/pkg/authorizer"

	"github.com/labstack/echo/v5"
)

// CasbinMiddleware enforces ABAC policies for every request.
// It expects AuthMiddleware to have already placed a UserAttr on the context.
type CasbinMiddleware struct {
	auth *authorizer.Authorizer
}

func NewCasbinMiddleware(auth *authorizer.Authorizer) *CasbinMiddleware {
	return &CasbinMiddleware{auth: auth}
}

func (m *CasbinMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			attr, ok := authorizer.FromContext(c.Request().Context())
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
			}

			obj := c.Request().URL.Path
			act := c.Request().Method

			allowed, err := m.auth.GetEnforcer().Enforce(attr, obj, act)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "authorization error"})
			}
			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
			}
			return next(c)
		}
	}
}
