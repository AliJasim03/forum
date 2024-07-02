package api

import (
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
	isLoggedIn := s.authenticateCookie(r)
	renderTemplate(w, "index", map[string]interface{}{
		"Title":      "Homepage",
		"isLoggedIn": isLoggedIn,
	})
}

func (s *server) registerPage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := s.authenticateCookie(r)
	renderTemplate(w, "register", map[string]interface{}{
		"Title":      "Register",
		"isLoggedIn": isLoggedIn,
	})
}

func (s *server) loginPage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := s.authenticateCookie(r)
	renderTemplate(w, "login", map[string]interface{}{
		"Title":      "Login",
		"isLoggedIn": isLoggedIn,
	})
}
