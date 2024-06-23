package forum

import (
	"forum/db"
	"net/http"
	"strings"
	"time"
)

func likePost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	DB := db.InitDB()
	defer db.CloseDB(DB)

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

	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	_, err = DB.Exec("INSERT INTO likes (user_id, post_id, is_like, created_at) VALUES(?,?,?,?)", userID, postID, true, time.Now())
	if err != nil {
		http.Error(res, "can't make like", http.StatusInternalServerError)
		return
	}
}

func likeComment(res http.ResponseWriter, req *http.Request) {
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

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	commentID := strings.Trim(req.URL.Path, "/cm")

	_, err = DB.Exec("INSERT INTO likes (user_id, commentt_id, is_like, created_at) VALUES(?,?,?,?)", userID, commentID, true, time.Now())
	if err != nil {
		http.Error(res, "can't make like", http.StatusInternalServerError)
		return
	}

}

func deslikePost(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	DB := db.InitDB()
	defer db.CloseDB(DB)

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

	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	_, err = DB.Exec("INSERT INTO likes (user_id, post_id, is_like, created_at) VALUES(?,?,?,?)", userID, postID, false, time.Now())
	if err != nil {
		http.Error(res, "can't make deslike", http.StatusInternalServerError)
		return
	}
}

func deslikeComment(res http.ResponseWriter, req *http.Request) {
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

	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		http.Error(res, "Session expired", http.StatusUnauthorized)
		return
	}

	commentID := strings.Trim(req.URL.Path, "/cm")

	_, err = DB.Exec("INSERT INTO likes (user_id, commentt_id, is_like, created_at) VALUES(?,?,?,?)", userID, commentID, false, time.Now())
	if err != nil {
		http.Error(res, "can't make dislike", http.StatusInternalServerError)
		return
	}
}
