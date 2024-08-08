package posts

type Topic struct {
	ID   int
	Name string
}

type Post struct {
	ID       int
	AuthorID int
	Title    string
	Content  string
	Topics   []*Topic
}

type Follower struct {
	UserID  int
	TopicID int
}
