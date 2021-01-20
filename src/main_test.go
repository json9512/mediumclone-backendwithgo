package main

import (
	"errors"
	"os"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/posts"
	"github.com/json9512/mediumclone-backendwithgo/src/tests"
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
		g.It("dbtool.CreateSamplePost should create a sample post in DB", func() {
			var post posts.PostModel
			posts.CreateTestSample(db)
			result := db.Where("author = ?", "test-author").Find(&post)

			g.Assert(result.Error).IsNil()
			g.Assert(post.Author).Eql("test-author")
			g.Assert(post.Comments).Eql(posts.JSONB{"comments-test": "testing 321"})
			g.Assert(post.Document).Eql(posts.JSONB{"testing": "test123"})
		})

		g.It("Delete the sample post created in DB", func() {
			var post posts.PostModel
			db.Where("author = ?", "test-author").Delete(&post)
			result := db.Where("author = ?", "test-author").Find(&post)
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

func runPostsTest(g *goblin.G, router *gin.Engine) {

}
