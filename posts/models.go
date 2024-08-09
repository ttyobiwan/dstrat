package posts

import "github.com/ttyobiwan/dstrat/users"

type Topic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	ID      int
	Title   string
	Content string
	Author  *users.User
	Topics  []*Topic
}

type Follower struct {
	UserID  int
	TopicID int
}
