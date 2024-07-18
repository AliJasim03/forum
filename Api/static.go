package api

import (
	"html/template"
	"net/http"
	"strconv"

	backend "forum/db"
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
	var posts []backend.Post
	backend.GetPosts(s.db, userID, &posts)
	Categories := backend.GetCategories(s.db)
	renderTemplate(w, "index", map[string]interface{}{
		"Title":      "Homepage",
		"isLoggedIn": isLoggedIn,
		"Posts":      posts,
		"Categories": Categories,
	})
}

func (s *server) filterCreatedPost(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, userID := s.authenticateCookie(r)
	var posts []backend.Post
	var filteredPosts []backend.Post
	backend.GetPosts(s.db, userID, &posts)
	// filter the posts that are created by the user
	for i := 0; i < len(posts); i++ {
		if posts[i].IsCreatedByUser {
			filteredPosts = append(filteredPosts, posts[i])
		}
	}
	Categories := backend.GetCategories(s.db)

	renderTemplate(w, "index", map[string]interface{}{
		"Title":      "Homepage",
		"isLoggedIn": isLoggedIn,
		"Posts":      filteredPosts,
		"Categories": Categories,
	})
}

func (s *server) filterLikedPost(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, userID := s.authenticateCookie(r)
	var posts []backend.Post
	backend.GetPosts(s.db, userID, &posts)
	var filteredPosts []backend.Post
	// filter posts that are created by the user
	for i := 0; i < len(posts); i++ {
		if posts[i].Like.IsLiked {
			filteredPosts = append(filteredPosts, posts[i])
		}
	}
	Categories := backend.GetCategories(s.db)

	renderTemplate(w, "index", map[string]interface{}{
		"Title":      "Homepage",
		"isLoggedIn": isLoggedIn,
		"Posts":      filteredPosts,
		"Categories": Categories,
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

func (s *server) postPage(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, userID := s.authenticateCookie(r)
	post := backend.Post{}
	//check if the post id is valid
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	// convert string to int
	id, err := strconv.Atoi(postID)
	if err != nil {
		// handle error
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
	}
	post.ID = id
	backend.GetPost(s.db, userID, &post)
	//check if post was not nil
	if post.Title == "" {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "postDetails", map[string]interface{}{
		"Title":      "Post",
		"isLoggedIn": isLoggedIn,
		"Post":       post,
	})
}
