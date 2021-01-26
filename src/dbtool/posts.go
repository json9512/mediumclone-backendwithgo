package dbtool

import (
	"time"

	"gorm.io/gorm"
)

// Post depicts "posts" table in db
type Post struct {
	ID        uint   `gorm:"primary_key"`
	Author    string `gorm:"column:author"`
	Document  JSONB  `gorm:"type:jsonb"`
	Comments  JSONB  `gorm:"type:jsonb"`
	Tags      string `gorm:"column:tags"`
	Likes     uint   `gorm:"column:likes"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CreateSamplePost creates a post sample in the database
func CreateSamplePost(db *Pool) {
	// test if post creation works
	doc := JSONB{"testing": "test123"}
	comments := JSONB{"comments-test": "testing 321"}

	post := Post{
		Author:   "test-author",
		Document: doc,
		Comments: comments,
		Likes:    0,
	}
	db.Insert(&post)
}
