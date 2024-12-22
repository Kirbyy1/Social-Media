package handlers

import (
	"fmt"
	"go-login-app/db"
	"go-login-app/models"
	"net/http"
	"time"
)

// CreatePostHandler handles the creation of a new post
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the content from the form
		content := r.FormValue("content")

		// Get the user from the session (assuming session management is in place)
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Create a new post
		post := models.Post{
			Username:  cookie.Value,
			Content:   content,
			CreatedAt: time.Now(),
		}

		// Save the post to the database
		if err := db.Create(&post); err != nil {
			http.Error(w, fmt.Sprintf("Error saving post: %v", err), http.StatusInternalServerError)
			return
		}

		// Redirect to the home page after posting
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

// GetPostsHandler fetches all posts for the feed
func GetPostsHandler() ([]models.Post, error) {
	return db.FindPosts()
}

// GetUserProfile fetches the user's profile
func GetUserProfile(username string) (models.User, error) {
	return db.FindUserByUsername(username)
}
