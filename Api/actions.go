package api

import (
	backend "forum/db"
	"net/http"
)

func (s *server) likeDislikePost(w http.ResponseWriter, r *http.Request) {
	//get the cookie to use token to get userID
	isLoggedIn, userID := s.authenticateCookie(r)
	if !isLoggedIn {
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
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
