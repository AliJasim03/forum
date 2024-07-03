package api

import (
	"encoding/json"
	backend "forum/db"
	"net/http"
)

func (s *server) likeDislikePost(w http.ResponseWriter, r *http.Request) {
	//get the cookie to use token to get userID
	isLoggedIn, userID := s.authenticateCookie(r)
	if !isLoggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var postID = r.FormValue("postID")
	if postID == "" {
		http.Error(w, "missing Post ID", http.StatusBadRequest)
		return
	}
	var action = r.FormValue("isLike")
	if action == "" {
		http.Error(w, "missing like or dislike", http.StatusBadRequest)
		return
	}
	var isLike = false
	if action == "like" {
		isLike = true
	} else if action == "dislike" {
		isLike = false
	}
	//save like to the database for the user
	ok := backend.LikeDislikePost(s.db, userID, postID, isLike)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Error(w, "can't make like", http.StatusInternalServerError)
}




func (s *server) createPost(res http.ResponseWriter, req *http.Request) {
	isLoggedIn, userID := s.authenticateCookie(req)
	if !isLoggedIn {
		http.Redirect(res, req, "/login", http.StatusUnauthorized)
		return
	}

	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}



	//get them values the body request json
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

	ok := backend.CreatePost(s.db, userID, post)
	if !ok {
		http.Error(res, "Failed to create post", http.StatusInternalServerError)
		return
	}
	//return message ok to the client
	res.WriteHeader(http.StatusOK)
}
