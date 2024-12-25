package main

import (
	"encoding/base64"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	tmpl = template.Must(template.ParseGlob("templates/*.html")) // Parses all .html files in the templates directory
	db   *gorm.DB
)

type User struct {
	ID                   uint
	Username             string `gorm:"unique"`
	Password             string
	ProfilePicture       string
	ProfilePictureData   []byte
	ProfilePictureBase64 string `gorm:"-"`
	Bio                  string
	Posts                []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	ID        uint
	Content   string
	UserID    uint
	User      User
	CreatedAt time.Time
}

func seedData() {
	user := User{
		Username: "testuser",
		Bio:      "This is a test user",
	}
	db.Create(&user)

	post := Post{
		Content: "This is a test post",
		UserID:  user.ID,
	}
	db.Create(&post)
}

//func fetchPosts() ([]Post, error) {
//	var posts []Post
//	result := db.Preload("User").Order("created_at desc").Find(&posts)
//	if result.Error != nil {
//		return nil, result.Error
//	}
//	return posts, nil
//}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("createPostHandler: Invalid method", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("createPostHandler: Form submitted")

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		log.Println("createPostHandler: Error parsing form", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")
	//log.Println("createPostHandler: Content received:", content)

	if content == "" {
		log.Println("createPostHandler: Empty content")
		http.Error(w, "Content cannot be empty", http.StatusBadRequest)
		return
	}

	// Retrieve the username from the cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("createPostHandler: Session cookie not found", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	username := cookie.Value
	log.Println("createPostHandler: Username retrieved from cookie:", username)

	// Fetch the user from the database
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		log.Println("createPostHandler: Error retrieving user from DB", result.Error)
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Create the post with the correct user ID
	post := Post{
		Content:   content,
		CreatedAt: time.Now(),
		UserID:    user.ID,
	}

	result = db.Create(&post)
	if result.Error != nil {
		log.Println("createPostHandler: Error creating post in DB", result.Error)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	log.Println("createPostHandler: Post created successfully for user:", username)

	// Redirect back to the feed or wherever appropriate
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the post ID from the form
	postID := r.FormValue("postID")

	// Check if the user is logged in
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve the user from the database
	var user User
	result := db.Where("username = ?", cookie.Value).First(&user)
	if result.Error != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if the post exists and belongs to the user
	var post Post
	result = db.Where("id = ? AND user_id = ?", postID, user.ID).First(&post)
	if result.Error != nil {
		http.Error(w, "Post not found or unauthorized", http.StatusNotFound)
		return
	}

	// Delete the post
	result = db.Delete(&post)
	if result.Error != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	// Redirect back to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Auto-migrate the User model
	db.AutoMigrate(&User{}, &Post{})

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already logged in by looking for a session cookie
	cookie, err := r.Cookie("session")
	if err == nil {
		// If the session cookie exists, we can assume the user is logged in
		var user User
		result := db.Where("username = ?", cookie.Value).First(&user)
		if result.Error == nil {
			// If a user is found, redirect them to the home page
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
	}

	// If no session or the user is not found, show the landing page
	err = tmpl.ExecuteTemplate(w, "landing.html", nil)
	if err != nil {
		log.Println("Error rendering landing template:", err)
	}
}

func getProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var user User
	result := db.Where("username = ?", cookie.Value).First(&user)
	if result.Error != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.ProfilePictureData != nil {
		user.ProfilePictureBase64 = "data:image/png;base64," + base64.StdEncoding.EncodeToString(user.ProfilePictureData)
		log.Println("Generated Base64 string for profile picture:", user.ProfilePictureBase64[:100]) // Debug log
	}
	log.Printf("Profile Picture Data Length: %d", len(user.ProfilePictureData))

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "upload.html", nil)
		return
	}

	log.Println("Parsing form data...")
	err := r.ParseMultipartForm(10 << 20) // 10MB size limit
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	log.Println("Retrieving file...")
	file, handler, err := r.FormFile("profilePicture")
	if err != nil {
		log.Println("Error retrieving file:", err)
		http.Error(w, "Unable to retrieve the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Printf("Received file: %s (%d bytes)", handler.Filename, handler.Size)

	// Get the user ID from the session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("User not logged in")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Find the user in the database
	var user User
	result := db.Where("username = ?", cookie.Value).First(&user)
	if result.Error != nil {
		log.Println("Error finding user:", result.Error)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Read the image file into a byte slice
	imageData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Unable to read the file", http.StatusInternalServerError)
		return
	}

	// Save the image filename and data in the database
	user.ProfilePicture = handler.Filename // Store the filename
	user.ProfilePictureData = imageData    // Store the image as binary data
	db.Save(&user)

	log.Println("File uploaded and user profile updated.")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "register.html", nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate input
	if username == "" || password == "" {
		tmpl.ExecuteTemplate(w, "register.html", "Username and password are required.")
		return
	}

	// Create a new user
	newUser := User{
		Username: username,
		Password: password, // Plain-text password (consider hashing in production)
	}

	// Attempt to save the user in the database
	result := db.Create(&newUser)
	if result.Error != nil {
		// Check if the error is due to a duplicate username
		if result.Error.Error() == "UNIQUE constraint failed: users.username" {
			tmpl.ExecuteTemplate(w, "register.html", "Username already exists. Please choose another.")
			return
		}

		log.Println("Error creating user:", result.Error)
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// Authenticate the user
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil || user.Password != password {
		tmpl.ExecuteTemplate(w, "login.html", "Invalid credentials. Please try again.")
		return
	}

	// Set a secure session (using HttpOnly and Secure flags for added security)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    username, // Store the username in the session
		Path:     "/",
		HttpOnly: true, // Only accessible by the server
		Secure:   true, // Only transmitted over HTTPS
		MaxAge:   3600, // Session expires in 1 hour
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the session cookie exists
	cookie, err := r.Cookie("session")
	if err != nil {
		log.Println("No session cookie found, redirecting to login.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch the user based on the session cookie (username)
	var user User
	result := db.Where("username = ?", cookie.Value).First(&user)
	if result.Error != nil {
		log.Printf("Error fetching user: %v", result.Error)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//log.Printf("User fetched: %+v", user)

	// Handle the post submission (POST request)
	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		post := Post{
			Content: content,
			UserID:  user.ID, // Associate post with the logged-in user
		}

		// Insert the post into the database
		result := db.Create(&post)
		if result.Error != nil {
			log.Printf("Error creating post: %v", result.Error)
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		// Redirect back to the home page to see the feed
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// Fetch posts with preloaded User data
	var posts []Post
	db.Preload("User").Find(&posts)
	for _, post := range posts {
		log.Printf("Post ID: %d, UserID: %d", post.ID, post.UserID)
	}
	result = db.Preload("User").Order("created_at desc").Find(&posts)
	if result.Error != nil {
		log.Printf("Error fetching posts: %v", result.Error)
	} else {
		log.Printf("Fetched posts: %+v", posts)
	}

	// Convert the image data to base64 encoding for the profile picture
	if user.ProfilePictureData != nil {
		user.ProfilePictureBase64 = base64.StdEncoding.EncodeToString(user.ProfilePictureData)
	}

	// Pass the user and posts to the template for rendering
	err = tmpl.ExecuteTemplate(w, "home.html", struct {
		User  User
		Posts []Post
	}{
		User:  user,
		Posts: posts,
	})
	if err != nil {
		log.Printf("Error rendering template: %v", err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Deletes the cookie
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Update the user's bio
func updateBioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the bio from the form submission
		bio := r.FormValue("bio")

		// Get the user from the session
		cookie, err := r.Cookie("session")
		if err != nil {
			log.Println("User not logged in")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		var user User
		result := db.Where("username = ?", cookie.Value).First(&user)
		if result.Error != nil {
			log.Println("Error finding user:", result.Error)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Update the bio in the database
		user.Bio = bio
		db.Save(&user)

		// Redirect back to home after updating bio
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func editProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve the user from the database
	var user User
	result := db.Where("username = ?", cookie.Value).First(&user)
	if result.Error != nil {
		log.Println("Error finding user:", result.Error)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// If the request is a GET, render the edit profile page
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "edit_profile.html", user)
		return
	}

	// If it's a POST request, handle the form submission
	err = r.ParseMultipartForm(10 << 20) // 10MB limit for form data
	if err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the new bio and profile picture from the form
	bio := r.FormValue("bio")
	file, _, err := r.FormFile("profilePicture")
	if err == nil {
		// If a new file is uploaded, handle it
		imageData, err := io.ReadAll(file)
		if err != nil {
			log.Println("Error reading file:", err)
			http.Error(w, "Unable to read the file", http.StatusInternalServerError)
			return
		}

		// Update the profile picture and bio in the database
		user.ProfilePictureData = imageData
		user.Bio = bio
	} else {
		// If no file is uploaded, just update the bio
		user.Bio = bio
	}

	db.Save(&user)

	// Redirect back to the home page
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func main() {
	initDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/edit-profile", editProfileHandler)
	http.HandleFunc("/update_bio", updateBioHandler)
	http.HandleFunc("/profile_picture", getProfilePictureHandler)
	http.HandleFunc("/create-post", createPostHandler)
	http.HandleFunc("/delete-post", deletePostHandler)

	log.Println("Server started on :8080")
	log.Println("Server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
