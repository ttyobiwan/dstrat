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

type PostHandler struct {
	store PostStore
}

func NewPostHandler(store PostStore) *PostHandler {
	return &PostHandler{store}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	// Get request data
	// TODO: Author should be taken from header or cookie
	data := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  int    `json:"author"`
		Topics  []int  `json:"topics"`
	}{}
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Detail string `json:"detail"`
		}{"Invalid request"})
	}

	// Create new post
	post, err := h.store.Create(data.Title, data.Content, data.Author, data.Topics)
	if err != nil {
		return fmt.Errorf("creating post: %v", err)
	}

	// TODO: Schedule task to notify followers
	// This could probably also require a transaction

	return c.JSON(http.StatusCreated, post)
}
