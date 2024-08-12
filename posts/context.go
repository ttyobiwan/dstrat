package posts

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type PostContext struct {
	echo.Context
	db *sql.DB
}

func (c *PostContext) DB() *sql.DB {
	return c.db
}
