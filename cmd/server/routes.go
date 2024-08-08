package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func addRoutes(e *echo.Echo) {
	// Internal
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Detail string `json:"detail"`
		}{"ok"})
	})
}
