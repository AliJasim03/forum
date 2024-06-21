package forum

import (
	"forum/db"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Username string
	Password string
}

type Post struct {
	PostID  int
	Title   string
	Content string
}

type Comment struct {
	PostID  int
	Content string
}

func Rigestrion(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//declare the db
	DB := db.InitDB()
	defer db.CloseDB(DB)

	email := req.FormValue("email")
	username := req.FormValue("username")
	password := req.FormValue("password")

	// check if all field used
	if email == "" || username == "" || password == "" {
		http.Error(res, "missing required fields", http.StatusBadRequest)
		return
	}

	// Check if the email is already used before
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? LIMIT 1)`
	err := DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(res, "email already registered", http.StatusConflict)
		return
	}

	//hash the pass to store it
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		return
	}

	//insert data
	_, err = DB.Exec("INSERT INTO users (email, username, password)VALUES(?, ?, ?)", email, username, hashPass)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		return
	}
}

func loggin(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//declare the db
	DB := db.InitDB()
	defer db.CloseDB(DB)

	username := req.FormValue("username")
	password := req.FormValue("password")

	// check if all field used
	if username == "" || password == "" {
		http.Error(res, "missing required fields", http.StatusBadRequest)
		return
	}

	var storedPass string
	var userID int
	err := DB.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &storedPass)
	if err != nil {
		http.Error(res, "invalild username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(password))
	if err != nil {
		http.Error(res, "invalild username or password", http.StatusUnauthorized)
		return
	}

	sessionToken, err := uuid.NewV4()
	if err != nil {
		http.Error(res, "failed to generate session token", http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)",
		userID, sessionToken.String(), time.Now().Add(10*time.Hour).Format("2000-08-11 17:01:09"))
	if err != nil {
		http.Error(res, "fail to store session in database", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    sessionToken.String(),
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(res, &cookie)

}

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

	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

}
