package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

var emptyCookie = &http.Cookie{
	Name:     "token", //l don't if l should change name or not because it same to the cookie of token
	Value:    "",
	Expires:  time.Unix(0, 0),
	HttpOnly: true,
	Path:     "/",
}

func (s *server) Registration(res http.ResponseWriter, req *http.Request) {

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
	err := s.db.QueryRow(query, email).Scan(&exists)
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
	err = s.db.QueryRow(query, username).Scan(&exists)
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
	_, err = s.db.Exec("INSERT INTO users (email, username, password)VALUES(?, ?, ?)", email, username, hashPass)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		return
	}

	//get the user id
	var userID int
	err = s.db.QueryRow("SELECT id FROM users WHERE username = ? and email = ?", username, email).Scan(&userID)
	if err != nil {
		http.Error(res, "Server error", http.StatusInternalServerError)
		return
	}
	cookie, err := s.generateCookie(fmt.Sprint(userID))
	if err != nil {
		http.Error(res, "fail to generate cookie", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &cookie)
}

func (s *server) Login(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	// check if all field used
	if email == "" || password == "" {
		http.Error(res, "missing required fields", http.StatusBadRequest)
		return
	}

	var storedPass string
	var userID int
	err := s.db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &storedPass)
	if err != nil {
		http.Error(res, "invalild username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(password))
	if err != nil {
		http.Error(res, "invalild username or password", http.StatusUnauthorized)
		return
	}

	cookie, err := s.generateCookie(fmt.Sprint(userID))
	if err != nil {
		http.Error(res, "fail to generate cookie", http.StatusInternalServerError)
		return
	}

	http.SetCookie(res, &cookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)

}

func (s *server) logout(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := req.Cookie("session_token")
	if err != nil {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	_, err = s.db.Exec("DELETE FROM sessions WHERE session_token = ?", sessionToken)
	if err != nil {
		http.Error(res, "fail to remove session from database", http.StatusInternalServerError)
		return
	}

	// put the empty cookie
	http.SetCookie(res, emptyCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func (s *server) generateCookie(userID string) (http.Cookie, error) {

	sessionToken, err := uuid.NewV4()
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed to generate session token: %w", err)
	}

	futureTime := time.Now().Add(10 * time.Hour)

	// Format the future time
	formattedTime := futureTime.Format("2006-01-02 15:04:05")

	s.db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)",
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

func (s *server) authenticateCookie(r *http.Request) bool {

	// extract token
	token, err := r.Cookie("token")
	if err != nil {
		return false
	}
	cookie := token
	sessionToken := cookie.Value
	var userID int
	var expiresAt time.Time

	//get the cookie to use token to get userID
	err = s.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&userID, &expiresAt)
	if err != nil || expiresAt.Before(time.Now()) {
		return false
	}

	// check if the session ended or not
	if expiresAt.Before(time.Now()) {
		return false
	}
	return true
}
