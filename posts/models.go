package posts

import "github.com/ttyobiwan/dstrat/users"

type Topic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	ID      int         `json:"id"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Author  *users.User `json:"author"`
	Topics  []*Topic    `json:"topics"`
}
