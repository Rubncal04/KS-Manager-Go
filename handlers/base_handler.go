package handlers

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Welcome() echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("user-info", c)

		if sess != nil {
			return c.Redirect(http.StatusSeeOther, "/user")
		}

		return c.Render(http.StatusOK, "login.html", nil)
	}
}
