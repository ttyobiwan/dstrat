package posts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TopicHandler struct{}

func NewTopicHandler() *TopicHandler {
	return &TopicHandler{}
}

func (h *TopicHandler) CreateTopic(c echo.Context) error {
	store := NewTopicDBStore(c.(*PostContext).DB())
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
	topic, err := store.GetByName(data.Name)
	if topic != nil {
		return c.JSON(http.StatusBadRequest, struct {
			Detail string `json:"detail"`
		}{"Such topic already exists"})
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return fmt.Errorf("getting topic: %v", err)
	}

	// Create new topic
	topic, err = store.Create(data.Name)
	if err != nil {
		return fmt.Errorf("creating topic: %v", err)
	}

	return c.JSON(http.StatusCreated, topic)
}

func (h *TopicHandler) FollowTopic(c echo.Context) error {
	return c.JSON(http.StatusCreated, nil)
}

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	store := NewPostDBStore(c.(*PostContext).DB())
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
	post, err := store.Create(data.Title, data.Content, data.Author, data.Topics)
	if err != nil {
		return fmt.Errorf("creating post: %v", err)
	}

	// TODO: Schedule task to notify followers
	// This could probably also require a transaction

	return c.JSON(http.StatusCreated, post)
}
