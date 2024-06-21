package forum

import (
	"forum/db"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Rigestrion (res http.ResponseWriter, req *http.Request) {
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

func Loggin (res http.ResponseWriter, req *http.Request) {
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

