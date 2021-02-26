package dbtool

import (
	"time"

	"gorm.io/gorm"
)

// Post depicts "posts" table in db
type Post struct {
	ID        uint   `gorm:"primary_key"`
	Author    string `gorm:"column:author"`
	Document  string `gorm:"column:document"`
	Comments  string `gorm:"column:comments"`
	Tags      string `gorm:"column:tags"`
	Likes     uint   `gorm:"column:likes"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CreatePost creates a new post with the given info
func (db *DB) CreatePost(doc, tags, author string) (*Post, error) {
	post := Post{Author: author, Document: doc, Tags: tags}
	query := db.Create(&post)

	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostByID returns the post with the provided id from DB
func (db *DB) GetPostByID(id uint) (*Post, error) {
	var post Post
	query := db.First(&post, "id = ?", id)
	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostByAuthor gets the post by the provided author
func (db *DB) GetPostByAuthor(author string) (*Post, error) {
	var post Post
	query := db.First(&post, "author = ?", author)

	if err := checkErr(query); err != nil {
		return nil, err
	}

	return &post, nil
}

// UpdatePost updates the post with the given data
func (db *DB) UpdatePost(newData interface{}) (*Post, error) {
	var updatedPost Post
	query := db.Model(&updatedPost).Updates(newData)
	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &updatedPost, nil
}
