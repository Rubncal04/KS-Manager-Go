package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authorization(requiredRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := Authentication(c)

			role := user.Role

			for _, requiredRole := range requiredRoles {
				if role == requiredRole {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "You don't have permission to access this route")
		}
	}
}
