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
	userHandler := users.NewUserHandler(users.NewUserDBStore(db))
	e.POST("/api/users", userHandler.CreateUser)

	// Posts
	topicHandler := posts.NewTopicHandler(posts.NewTopicDBStore(db))
	e.POST("/api/topics", topicHandler.CreateTopic)
}
