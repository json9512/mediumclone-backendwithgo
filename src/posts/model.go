package posts

import (
	"github.com/jinzhu/gorm"
	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

// Post depicts "posts" table in db
type Post struct {
	gorm.Model
	Author   string       `gorm:"column:author"`
	Document config.JSONB `gorm:"type:jsonb"`
	Comments config.JSONB `gorm:"type:jsonb"`
	Likes    uint         `gorm:"column:likes"`
}

// CreateTestSample creates a post sample in the database
func CreateTestSample(db *gorm.DB) {
	// test if post creation works
	doc := config.JSONB{"testing": "test123"}
	comments := config.JSONB{"comments-test": "testing 321"}

	post := Post{
		Author:   "test-author",
		Document: doc,
		Comments: comments,
		Likes:    0,
	}
	db.Create(&post)
}
