package api

import (
	"encoding/json"
	"net/http"

	backend "forum/db"
)

type PostJson struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

type LikeDisJson struct {
	ID     string `json:"ID"`
	IsLike string `json:"isLike"`
}

func (s *server) likeDislikePost(w http.ResponseWriter, r *http.Request) {
	// get the cookie to use token to get userID
	isLoggedIn, userID := s.authenticateCookie(r)

	var LikeDis LikeDisJson

	err := json.NewDecoder(r.Body).Decode(&LikeDis)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if !isLoggedIn {
		http.Error(w, "Please log in to continue", http.StatusBadRequest)
		return
	}

	action := LikeDis.IsLike
	if action == "" || LikeDis.ID == "" {
		http.Error(w, "missing like or dislike", http.StatusBadRequest)
		return
	}
	isLike := false
	if LikeDis.IsLike == "like" {
		isLike = true
	} else if LikeDis.IsLike == "dislike" {
		isLike = false
	}
	// save like to the database for the user
	ok := backend.LikeDislikePost(s.db, userID, LikeDis.ID, isLike)
	if ok {

		isLiked := backend.KnowPostLike(s.db, userID, LikeDis.ID)
		// return data to the client that the like is success
		w.Header().Set("Content-Type", "application/json")
		// return isliked
		json.NewEncoder(w).Encode(isLiked)
		return
	}
	http.Error(w, "can't make like", http.StatusInternalServerError)
}

func (s *server) createPost(res http.ResponseWriter, req *http.Request) {
	isLoggedIn, userID := s.authenticateCookie(req)
	if !isLoggedIn {
		http.Error(res, "Please log in to continue", http.StatusBadRequest)
		return
	}

	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get them values the body request json
	var post backend.PostJson
	err := json.NewDecoder(req.Body).Decode(&post)
	if err != nil {
		http.Error(res, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Check if required fields are provided
	if post.Title == "" || post.Content == "" {
		http.Error(res, "Title & content are required", http.StatusBadRequest)
		return
	}

	postId := backend.CreatePost(s.db, userID, post)
	if postId == -1 {
		http.Error(res, "Failed to create post", http.StatusInternalServerError)
		return
	}
	// return message ok to the client
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(postId)
}

func (s *server) createComment(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	isLoggedIn, userID := s.authenticateCookie(req)
	if !isLoggedIn {
		http.Error(res, "Please log in to continue", http.StatusBadRequest)
		return
	}

	var comment backend.CommentJson
	err := json.NewDecoder(req.Body).Decode(&comment)
	if err != nil {
		http.Error(res, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	if comment.Comment == "" || comment.PostID == "" {
		http.Error(res, "PostID & comment are required", http.StatusBadRequest)
		return
	}
	ok, retunedComment := backend.CreateComment(s.db, userID, comment)
	if !ok {
		http.Error(res, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// return the comment to the client
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(retunedComment)
}

func (s *server) likeDislikeComment(w http.ResponseWriter, r *http.Request) {
	// get the cookie to use token to get userID
	isLoggedIn, userID := s.authenticateCookie(r)

	var LikeDis LikeDisJson

	err := json.NewDecoder(r.Body).Decode(&LikeDis)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if !isLoggedIn {
		http.Error(w, "Please log in to continue", http.StatusBadRequest)
		return
	}

	if LikeDis.IsLike == "" || LikeDis.ID == "" {
		http.Error(w, "missing like or dislike", http.StatusBadRequest)
		return
	}
	isLike := false
	if LikeDis.IsLike == "like" {
		isLike = true
	} else if LikeDis.IsLike == "dislike" {
		isLike = false
	}
	// save like to the database for the user
	ok := backend.LikeDislikeComment(s.db, userID, LikeDis.ID, isLike)
	if ok {
		isLiked := backend.KnowCommentLike(s.db, userID, LikeDis.ID)
		// return data to the client that the like is success
		w.Header().Set("Content-Type", "application/json")
		// return isliked
		json.NewEncoder(w).Encode(isLiked)
		return
	}
	http.Error(w, "can't make like", http.StatusInternalServerError)
}

func (s *server) getPostLikesAndDislikesCount(w http.ResponseWriter, r *http.Request) {
	var ID backend.IDJson

	err := json.NewDecoder(r.Body).Decode(&ID)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if ID.ID == "" {
		http.Error(w, "missing postID", http.StatusBadRequest)
		return
	}
	// get the count of likes and dislikes
	likes, dislike := backend.GetPostLikesAndDislikesCount(s.db, ID.ID)
	// return the count to the client
	w.Header().Set("Content-Type", "application/json")
	// build json
	counts := map[string]int{
		"likes":    likes,
		"dislikes": dislike,
	}
	json.NewEncoder(w).Encode(counts)
}

func (s *server) getCommentLikesAndDislikesCount(w http.ResponseWriter, r *http.Request) {
	var ID backend.IDJson

	err := json.NewDecoder(r.Body).Decode(&ID)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if ID.ID == "" {
		http.Error(w, "missing commentID", http.StatusBadRequest)
		return
	}
	// get the count of likes and dislikes
	likes, dislike := backend.GetCommentLikesAndDislikesCount(s.db, ID.ID)
	// return the count to the client
	w.Header().Set("Content-Type", "application/json")
	// build json
	counts := map[string]int{
		"likes":    likes,
		"dislikes": dislike,
	}
	json.NewEncoder(w).Encode(counts)
}
