package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// note: make sure to check if the user session expired or not
type PostJson struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Gategories []string `json:"categories"`
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
	var post PostJson
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

	//create the post
	_, err = s.db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, post.Title, post.Content)
	if err != nil {
		http.Error(res, "Failed to create post", http.StatusInternalServerError)
		return
	}

	var postID int
	err = s.db.QueryRow("SELECT id FROM posts WHERE user_id = ? AND title = ? ORDER BY created_at DESC LIMIT 1", userID, post.Title).Scan(&postID)
	if err != nil {
		http.Error(res, "Failed to retrieve post ID", http.StatusInternalServerError)
		return
	}

	//add the category
	for _, ct := range post.Gategories {
		if ct != "" { // check if not empty
			var categoryID int
			err = s.db.QueryRow("SELECT id FROM categories WHERE name = ?", ct).Scan(&categoryID)
			if err != nil {
				http.Error(res, "Invalid category", http.StatusBadRequest)
				return
			}

			_, err = s.db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
			if err != nil {
				http.Error(res, "Failed to associate category with post", http.StatusInternalServerError)
				return
			}
		}
	}
	//return message ok to the client
	res.WriteHeader(http.StatusOK)
}

func (s *server) createComment(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//get the cookie to use token to get userID
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	//get the user_id of who will comment
	err = s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	content := req.FormValue("content")
	if content == "" {
		http.Error(res, "the content must be fill", http.StatusBadRequest)
		return
	}

	//here should use strings to get specfic prefix of URL
	postID := strings.Trim(req.URL.Path, "/lg")

	_, err = s.db.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?)", postID, userID, content, time.Now())
	if err != nil {
		http.Error(res, "fail to insert data to comment table", http.StatusInternalServerError)
	}
}

func (s *server) likeDislikePost_sorryIdidntsawit(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//here should use strings to get specfic prefix of URL
	postID := strings.Trim(req.URL.Path, "/lg")

	//get the cookie to use token to get userID
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	err = s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	_, err = s.db.Exec("INSERT INTO likes (user_id, post_id, is_like, created_at) VALUES(?,?,?,?)", userID, postID, true, time.Now())
	if err != nil {
		http.Error(res, "can't make like", http.StatusInternalServerError)
		return
	}
}

func (s *server) likeComment(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//get the cookie to use token to get userID
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	//get the user_id of who will comment
	err = s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	commentID := strings.Trim(req.URL.Path, "/cm")

	_, err = s.db.Exec("INSERT INTO likes (user_id, commentt_id, is_like, created_at) VALUES(?,?,?,?)", userID, commentID, true, time.Now())
	if err != nil {
		http.Error(res, "can't make like", http.StatusInternalServerError)
		return
	}

}

func (s *server) dislikePost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//here should use strings to get specfic prefix of URL
	postID := strings.Trim(req.URL.Path, "/lg")

	//get the cookie to use token to get userID
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	err = s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	_, err = s.db.Exec("INSERT INTO likes (user_id, post_id, is_like, created_at) VALUES(?,?,?,?)", userID, postID, false, time.Now())
	if err != nil {
		http.Error(res, "can't make deslike", http.StatusInternalServerError)
		return
	}
}

func (s *server) dislikeComment(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//get the cookie to use token to get userID
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	//get the user_id of who will comment
	err = s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	commentID := strings.Trim(req.URL.Path, "/cm")

	_, err = s.db.Exec("INSERT INTO likes (user_id, commentt_id, is_like, created_at) VALUES(?,?,?,?)", userID, commentID, false, time.Now())
	if err != nil {
		http.Error(res, "can't make dislike", http.StatusInternalServerError)
		return
	}
}
