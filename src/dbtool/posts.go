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

func (db *DB) CreatePost(doc, tags, author string) (*Post, error) {
	post := Post{Author: author, Document: doc, Tags: tags}
	query := db.Create(&post)

	if err := checkErr(query); err != nil {
		return nil, err
	}
	return &post, nil
}

func (db *DB) GetPostByAuthor(author string) (*Post, error) {
	var post Post
	query := db.First(&post, "author = ?", author)

	if err := checkErr(query); err != nil {
		return nil, err
	}

	return &post, nil
}
