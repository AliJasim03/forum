package db

type Post struct {
	ID              int
	Title           string
	Content         string
	CreatedBy       string
	CreatedOn       string
	Categories      []string
	Like            Like
	Comments        []Comments
	IsCreatedByUser bool
}
type Comments struct {
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
