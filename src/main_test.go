package main

import (
	"errors"
	"os"
	"testing"

	"github.com/franela/goblin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/posts"
	"github.com/json9512/mediumclone-backendwithgo/src/tests"
)

func Test(t *testing.T) {
	config.ReadVariablesFromFile(".env")
	// Setup test_db for local use only
	db := dbtool.Init()
	dbtool.Migrate(db)

	// Setup router
	router := SetupRouter("test", db)

	// create goblin
	g := goblin.Goblin(t)
	tests.RunPostsTests(g, router)
	tests.RunUsersTests(g, router, db)

	// NOTE for future me: For testing only,
	// - all the test cases in DB should move to each endpoint
	// - CRUD on endpoints should test the interaction with DB
	// - no separate DB test is required
	g.Describe("DB test", func() {
		g.It("posts.CreateTestSample should create a sample post in DB", func() {
			var post posts.Post
			posts.CreateTestSample(db)
			result := db.Where("author = ?", "test-author").Find(&post)

			g.Assert(result.Error).IsNil()
			g.Assert(post.Author).Eql("test-author")
			g.Assert(post.Comments).Eql(config.JSONB{"comments-test": "testing 321"})
			g.Assert(post.Document).Eql(config.JSONB{"testing": "test123"})
		})

		g.It("Delete the sample post created in DB", func() {
			var post posts.Post
			db.Where("id = ?", "1").Delete(&post)
			db.Unscoped().Delete(&post)
			result := db.Where("id = ?", "1").Find(&post)
			g.Assert(result.Error).Eql(errors.New("record not found"))
		})
	})

	tests.RunAuthTests(g, router, db)

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

	// Drop the users table
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE posts")

	// Note: should separate the test db and production db
}
