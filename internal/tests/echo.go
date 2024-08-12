package tests

import (
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
)

func MakeRequest(method, target, path, data string) (*httptest.ResponseRecorder, echo.Context) {
	e := echo.New()
	req := httptest.NewRequest(method, target, strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if path != "" {
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(path)
	}
	return rec, c
}
