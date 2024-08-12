package posts

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func GetRoutes(e *echo.Echo, db *sql.DB) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&PostContext{c, db})
		}
	})
	topicHandler := NewTopicHandler()
	e.POST("/api/topics", topicHandler.CreateTopic)
	e.POST("/api/topics/:id/follow", topicHandler.FollowTopic)
	postHandler := NewPostHandler()
	e.POST("/api/posts", postHandler.CreatePost)
}
