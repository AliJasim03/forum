package server

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
	mux.Handle("/front-end/static/", http.StripPrefix("/front-end/static/", http.FileServer(http.Dir("static"))))

	// Define routes
	mux.HandleFunc("/", indexHandler)

	fmt.Println("Server is running on http://localhost:8080/")
	//open in browser
	open("http://localhost:8080/")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpl = "templates/" + tmpl + ".html"
	t, err := template.ParseFiles(tmpl, "templates/layout.html")
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
	renderTemplate(w, "index", map[string]interface{}{
		"Title": "Homepage",
	})
}
