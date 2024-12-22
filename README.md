# Social Media Feed Application

This repository contains a simple web application where users can post content, upload profile pictures, and view posts in a feed.

## Features

- User registration and login.
- Create, view, and manage posts.
- Profile picture upload and display.
- Edit user profile details, including bio.

## Preview

### User Feed

![User Feed Preview](static/images/feed_example.png)

### Edit Profile Page

![Edit Profile Preview](static/images/edit_profile_example.png)

## Usage

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Run the application:
   ```bash
   go run main.goa
   ```

3. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## Folder Structure

- **templates/**: Contains all HTML templates.
- **static/**: Holds static files (CSS, images, JavaScript).
- **main.go**: The main Go application file.

## Acknowledgements

This project was created as a demonstration of a basic social media feed with Go, GORM, and SQLite.
