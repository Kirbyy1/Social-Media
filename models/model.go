package models

import "time"

// User model
type User struct {
	ID                   uint   `gorm:"primaryKey"`
	Username             string `gorm:"unique;not null"`
	Bio                  string
	ProfilePictureData   []byte // Profile picture data in binary form
	ProfilePictureBase64 string
}

// Post model
type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:current_timestamp"`
}
