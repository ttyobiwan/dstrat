package posts

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/client"
)

type TopicHandler struct{}

func NewTopicHandler() *TopicHandler {
	return &TopicHandler{}
}

func (h *TopicHandler) CreateTopic(c echo.Context) error {
	store := c.(*PostContext).TopicStore()
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
	store := c.(*PostContext).TopicStore()

	topic_id := c.Param("id")
	// TODO: This should rather by set in some middleware but whatever
	user_id := c.Request().Header.Get("user")

	err := store.ToggleFollow(topic_id, user_id)
	if err != nil {
		return fmt.Errorf("toggling follow: %v", err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

type PostHandler struct{}

func NewPostHandler() *PostHandler {
	return &PostHandler{}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	pctx := c.(*PostContext)

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
	store := pctx.PostStore()
	post, err := store.Create(data.Title, data.Content, data.Author, data.Topics)
	if err != nil {
		return fmt.Errorf("creating post: %v", err)
	}

	// Send post to the followers
	// TODO: Add tx and think about the result
	tc := pctx.TaskClient()
	tc.Execute(
		c.Request().Context(),
		SendPostToTopicFollowers,
		client.StartWorkflowOptions{TaskQueue: pctx.TaskQueue()},
		post.ID,
		data.Topics,
	)

	return c.JSON(http.StatusCreated, post)
}
