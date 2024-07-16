package db

type Post struct {
	ID              int
	Title           string
	Content         string
	CreatedBy       string
	CreatedOn       string
	Categories      []string
	Like            Like
	Comments        []Comment
	IsCreatedByUser bool
}

type Comment struct {
	ID              int
	PostID          int
	CreatedBy       string
	Content         string
	CreatedOn       string
	Like            Like
	IsCreatedByUser bool
}

type Like struct {
	CountLikes    int
	CountDislikes int
	IsLiked       bool
	IsDisliked    bool
}

type Category struct {
	ID   int
	Name string
}

type PostJson struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

type CommentJson struct {
	PostID  string `json:"PostID"`
	Comment string `json:"comment"`
}

type IDJson struct {
	ID string `json:"ID"`
}
