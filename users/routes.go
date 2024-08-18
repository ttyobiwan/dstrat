package users

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func GetRoutes(e *echo.Echo, db *sql.DB) {
	g := e.Group("/api", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&UserContext{c, db})
		}
	})
	userHandler := NewUserHandler()
	g.POST("/users", userHandler.CreateUser)
}
