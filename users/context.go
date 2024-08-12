package users

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

type UserContext struct {
	echo.Context
	db *sql.DB
}

func (c *UserContext) DB() *sql.DB {
	return c.db
}
