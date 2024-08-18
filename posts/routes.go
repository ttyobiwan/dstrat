package posts

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/ttyobiwan/dstrat/internal/temporal"
)

func GetRoutes(e *echo.Echo, db *sql.DB, tc *temporal.Client) {
	g := e.Group("/api", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&PostContext{c, db, tc, temporal.TemporalQueuePosts})
		}
	})
	topicHandler := NewTopicHandler()
	g.POST("/topics", topicHandler.CreateTopic)
	g.POST("/topics/:id/follow", topicHandler.FollowTopic)
	postHandler := NewPostHandler()
	g.POST("/posts", postHandler.CreatePost)
}
