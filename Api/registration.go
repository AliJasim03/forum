package forum

import (
	"fmt"
	"forum/db"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var dcookie = &http.Cookie{
	Name:     "token", //l don't if l should change name or not because it same to the cookie of token
	Value:    "",
	Expires:  time.Unix(0, 0),
	HttpOnly: true,
	Path:     "/",
}

func Rigestrion(res http.ResponseWriter, req *http.Request) {
	DB := db.InitDB()
	defer db.CloseDB(DB)
	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// Check if the name is already used before
	query = `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? LIMIT 1)`
	err = DB.QueryRow(query, username).Scan(&exists)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
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

	//get the user id
	var userID int
	err = DB.QueryRow("SELECT id FROM users WHERE username = ? and email = ?", username, email).Scan(&userID)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		return
	}
	cookie, err := generateCookie(fmt.Sprint(userID))
	if err != nil {
		http.Error(res, "fail to generate cookie", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &cookie)
}

func Login(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	DB := db.InitDB()
	//declare the db
	defer db.CloseDB(DB)

	email := req.FormValue("email")
	password := req.FormValue("password")

	// check if all field used
	if email == "" || password == "" {
		http.Error(res, "missing required fields", http.StatusBadRequest)
		return
	}

	var storedPass string
	var userID int
	err := DB.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storedPass)
	if err != nil {
		http.Error(res, "invalild username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(password))
	if err != nil {
		http.Error(res, "invalild username or password", http.StatusUnauthorized)
		return
	}

	cookie, err := generateCookie(fmt.Sprint(userID))
	if err != nil {
		http.Error(res, "fail to generate cookie", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &cookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)

}

func Logout(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//destroy the cookie
	http.SetCookie(res, dcookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func generateCookie(userID string) (http.Cookie, error) {

	DB := db.InitDB()
	//declare the db
	defer db.CloseDB(DB)

	sessionToken, err := uuid.NewV4()
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed to generate session token: %w", err)
	}

	futureTime := time.Now().Add(10 * time.Hour)

	// Format the future time
	formattedTime := futureTime.Format("2006-01-02 15:04:05")

	DB.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)",
		userID, sessionToken.String(), formattedTime)

	cookie := http.Cookie{
		Name:     "token",
		Value:    sessionToken.String(),
		Expires:  time.Now().Add(10 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	return cookie, nil
}

func authenticateCookie(r *http.Request) bool {
	//extracat tocken
	tocken, err := r.Cookie("token")
	if err != nil {
		return false
	}
	//get the cookie to use token to get userID
	DB := db.InitDB()
	defer db.CloseDB(DB)
	cookie := tocken
	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	err = DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		return false
	}
	//check if the session already end or not
	if expiresAt.Before(time.Now()) {
		return false
	}
	return true
}
