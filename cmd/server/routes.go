package main

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ttyobiwan/dstrat/posts"
	"github.com/ttyobiwan/dstrat/users"
)

func addRoutes(e *echo.Echo, db *sql.DB) {
	// Internal
	e.GET("/api/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Detail string `json:"detail"`
		}{"ok"})
	})

	// Users
	users.GetRoutes(e, db)

	// Posts
	posts.GetRoutes(e, db)
}
