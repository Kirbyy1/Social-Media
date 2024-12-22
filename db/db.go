package db

import (
	"go-login-app/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// Init initializes the database connection
func Init() {
	var err error
	// Open a SQLite database (you can use any database here)
	DB, err = gorm.Open(sqlite.Open("gypsi_network.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// AutoMigrate: Automatically create tables for models if they don't exist
	DB.AutoMigrate(&models.User{}, &models.Post{})
}

// Create is a wrapper for DB create operation
func Create(value interface{}) error {
	result := DB.Create(value)
	return result.Error
}

// FindPosts retrieves all posts ordered by creation time in descending order
func FindPosts() ([]models.Post, error) {
	var posts []models.Post
	result := DB.Order("created_at desc").Find(&posts) // Correct use of `Order`
	return posts, result.Error
}

// FindUserByUsername finds a user by their username
func FindUserByUsername(username string) (models.User, error) {
	var user models.User
	result := DB.Where("username = ?", username).First(&user) // Correct use of `Where`
	return user, result.Error
}
