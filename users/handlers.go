package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	store UserStore
}

func NewUserHandler(store UserStore) *UserHandler {
	return &UserHandler{store}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	// Get request data
	data := struct {
		Username string `json:"username"`
	}{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Detail string `json:"detail"`
		}{"Invalid request"})
	}

	// Check username uniqueness
	user, err := h.store.GetByUsername(data.Username)
	if user != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Detail string `json:"detail"`
		}{"Username already taken"})
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("getting user: %v", err)
	}

	// Create new user
	user, err = h.store.Create(data.Username)
	if err != nil {
		return fmt.Errorf("creating user: %v", err)
	}

	return c.JSON(http.StatusCreated, user)
}
