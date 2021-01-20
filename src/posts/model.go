package posts

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// PostModel table for Post in db
type PostModel struct {
	gorm.Model
	Author   string `gorm:"column:author"`
	Document JSONB  `gorm:"type:jsonb"`
	Comments JSONB  `gorm:"type:jsonb"`
	Likes    uint   `gorm:"column:likes"`
}

// JSONB type def for database
type JSONB map[string]interface{}

// Value for gorm to read the JSONB data
func (j JSONB) Value() (driver.Value, error) {
	valString, err := json.Marshal(j)
	return string(valString), err
}

// Scan for gorm to scan the JSONB data
func (j *JSONB) Scan(v interface{}) error {
	err := json.Unmarshal(v.([]byte), &j)
	if err != nil {
		return err
	}
	return nil
}

// CreateTestSample creates test post sample in database
func CreateTestSample(db *gorm.DB) {
	// test if post creation works
	doc := JSONB{"testing": "test123"}
	comments := JSONB{"comments-test": "testing 321"}

	post := PostModel{
		Author:   "test-author",
		Document: doc,
		Comments: comments,
		Likes:    0,
	}
	db.Create(&post)
}
