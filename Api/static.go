package api

import (
	backend "forum/db"
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = "front-end/templates/" + tmpl + ".html"
	t, err := template.ParseFiles(tmpl, "front-end/templates/layout.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, userID := s.authenticateCookie(r)
	renderTemplate(w, "index", map[string]interface{}{
		"Title":      "Homepage",
		"isLoggedIn": isLoggedIn,
		"Posts":      s.getPosts(userID, ""),
	})
}

func (s *server) registerPage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, _ := s.authenticateCookie(r)
	renderTemplate(w, "register", map[string]interface{}{
		"Title":      "Register",
		"isLoggedIn": isLoggedIn,
	})
}

func (s *server) loginPage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, _ := s.authenticateCookie(r)
	renderTemplate(w, "login", map[string]interface{}{
		"Title":      "Login",
		"isLoggedIn": isLoggedIn,
	})
}
func (s *server) createPostPage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, _ := s.authenticateCookie(r)
	categories := backend.GetCategories(s.db)
	renderTemplate(w, "createPost", map[string]interface{}{
		"Title":      "Create Post",
		"isLoggedIn": isLoggedIn,
		"Categories": categories,
	})
}
