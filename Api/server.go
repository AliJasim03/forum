package forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func ServerInit() {
	// Create a new ServeMux
	mux := http.NewServeMux()
	// Serve static files (CSS, JavaScript, etc.)
	mux.Handle("/front-end/static/", http.StripPrefix("/front-end/static/", http.FileServer(http.Dir("./front-end/static"))))

	// Define routes
	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/register", registerPage)

	mux.HandleFunc("/registerAction", Rigestrion)

	mux.HandleFunc("/login", loginPage)

	mux.HandleFunc("/loginAction", Login)

	mux.HandleFunc("/logout", Logout)

	fmt.Println("Server is running on http://localhost:8080/")
	//open in browser
	open("http://localhost:8080/")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	isLooged := authenticateCookie(r)
	RenderTemplate(w, "index", map[string]interface{}{
		"Title":      "Homepage",
		"isLoggedIn": isLooged,
	})
}

func registerPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "register", map[string]interface{}{
		"Title": "Register",
	})
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "login", map[string]interface{}{
		"Title": "Login",
	})
}
