package posts

import (
	"database/sql"

	"github.com/labstack/echo/v4"
	"github.com/ttyobiwan/dstrat/internal/temporal"
	"go.temporal.io/sdk/client"
)

type PostContext struct {
	echo.Context
	db *sql.DB
	tc *temporal.Client
	tq temporal.TemporalQueue
}

func (c *PostContext) DB() *sql.DB {
	return c.db
}

func (c *PostContext) TaskQueue() string {
	return string(c.tq)
}

func (c *PostContext) TaskClient() TaskClient[client.StartWorkflowOptions, client.WorkflowRun] {
	return c.tc
}

func (c *PostContext) TopicStore() TopicStore {
	return NewTopicDBStore(c.DB())
}

func (c *PostContext) PostStore() PostStore {
	return NewPostDBStore(c.DB())
}
