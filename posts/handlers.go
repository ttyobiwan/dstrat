package posts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TopicHandler struct {
	store TopicStore
}

func NewTopicHandler(store TopicStore) *TopicHandler {
	return &TopicHandler{store}
}

func (h *TopicHandler) CreateTopic(c echo.Context) error {
	// Get request data
	data := struct {
		Name string `json:"name"`
	}{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Detail string `json:"detail"`
		}{"Invalid request"})
	}

	// Check name uniqueness
	topic, err := h.store.GetByName(data.Name)
	if topic != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Detail string `json:"detail"`
		}{"Such topic already exists"})
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("getting topic: %v", err)
	}

	// Create new topic
	topic, err = h.store.Create(data.Name)
	if err != nil {
		return fmt.Errorf("creating topic: %v", err)
	}

	return c.JSON(http.StatusCreated, topic)
}
