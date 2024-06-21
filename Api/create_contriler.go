package forum

import (
	"forum/db"
	"net/http"
	"time"
)

func createPost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//declare the db
	DB := db.InitDB()
	defer db.CloseDB(DB)

	//get the cookie to use token to get userID
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	title := req.FormValue("title")
	content := req.FormValue("content")
	catType := req.Form["categories"]

	// Check if required fields are provided
	if title == "" || content == "" {
		http.Error(res, "Title & content are required", http.StatusBadRequest)
		return
	}

	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	//get the user_id of who will post
	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//create the post
	_, err = DB.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", userID, title, content)
	if err != nil {
		http.Error(res, "Failed to create post", http.StatusInternalServerError)
		return
	}

	var postID int
	err = DB.QueryRow("SELECT id FROM posts WHERE user_id = ? AND title = ? ORDER BY created_at DESC LIMIT 1", userID, title).Scan(&postID)
	if err != nil {
		http.Error(res, "Failed to retrieve post ID", http.StatusInternalServerError)
		return
	}

	//add the category
	for _, ct := range catType {
		if ct != "" { // check if not empty
			var categoryID int
			err = DB.QueryRow("SELECT id FROM categories WHERE name = ?", ct).Scan(&categoryID)
			if err != nil {
				http.Error(res, "Invalid category", http.StatusBadRequest)
				return
			}

			_, err = DB.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
			if err != nil {
				http.Error(res, "Failed to associate category with post", http.StatusInternalServerError)
				return
			}
		}
	}

}

func createComment(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//declare the db
	DB := db.InitDB()
	defer db.CloseDB(DB)

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
	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	content := req.FormValue("content")
	if content == "" {
		http.Error(res, "the content must be fill", http.StatusBadRequest)
		return
	}

}
