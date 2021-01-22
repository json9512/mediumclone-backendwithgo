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
	"github.com/json9512/mediumclone-backendwithgo/src/users"
)

func Test(t *testing.T) {
	// Setup router
	router := SetupRouter("test")

	// Setup test_db for local use only
	db := dbtool.Init()
	dbtool.Migrate(db)

	// create goblin
	g := goblin.Goblin(t)
	tests.RunPostsTests(g, router)
	tests.RunUsersTests(g, router)

	g.Describe("DB test", func() {
		g.It("posts.CreateTestSample should create a sample post in DB", func() {
			var post posts.PostModel
			posts.CreateTestSample(db)
			result := db.Where("author = ?", "test-author").Find(&post)

			g.Assert(result.Error).IsNil()
			g.Assert(post.Author).Eql("test-author")
			g.Assert(post.Comments).Eql(config.JSONB{"comments-test": "testing 321"})
			g.Assert(post.Document).Eql(config.JSONB{"testing": "test123"})
		})

		g.It("Delete the sample post created in DB", func() {
			var post posts.PostModel
			db.Where("author = ?", "test-author").Delete(&post)
			result := db.Where("author = ?", "test-author").Find(&post)
			g.Assert(result.Error).Eql(errors.New("record not found"))
		})

		g.It("users.CreateTestSample should create a sample user in DB", func() {
			var user users.UserModel
			users.CreateTestSample(db)
			result := db.Where("username = ?", "test-user").Find(&user)

			g.Assert(result.Error).IsNil()
			g.Assert(user.Email).Eql("test@test.com")
			g.Assert(user.Username).Eql("test-user")
		})

		g.It("Delete the sample user created in DB", func() {
			var user users.UserModel
			db.Where("username = ?", "test-user").Delete(&user)
			result := db.Where("username = ?", "test-user").Find(&user)
			g.Assert(result.Error).Eql(errors.New("record not found"))
		})
	})

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

}
