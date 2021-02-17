package main

import (
	"errors"
	"os"
	"testing"

	"github.com/franela/goblin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
	"github.com/json9512/mediumclone-backendwithgo/src/tests"
)

func Test(t *testing.T) {
	config.ReadVariablesFromFile(".env")
	// Setup test_db for local use only
	pool := dbtool.Init()

	// NOTE: need to copy existing data to temp table
	// Drop the tables
	pool.Exec("DROP TABLE users")
	pool.Exec("DROP TABLE posts")

	dbtool.Migrate(pool)
	router := SetupRouter("test", pool)
	g := goblin.Goblin(t)

	toolBox := tests.TestToolbox{
		Goblin: g,
		Router: router,
		DB:     pool,
	}
	tests.RunPostsTests(&toolBox)
	tests.RunUsersTests(&toolBox)

	// NOTE for future me: For testing only,
	// - all the test cases in DB should move to each endpoint
	// - CRUD on endpoints should test the interaction with DB
	// - no separate DB test is required
	g.Describe("DB test", func() {
		g.It("CreateSamplePost should create a sample post in DB", func() {
			var post dbtool.Post
			dbtool.CreateSamplePost(pool)
			dbErr := pool.Query(
				&post,
				map[string]interface{}{"author": "test-author"},
			)

			g.Assert(dbErr).IsNil()
			g.Assert(post.Author).Eql("test-author")
			g.Assert(post.Comments).Eql(dbtool.JSONB{"comments-test": "testing 321"})
			g.Assert(post.Document).Eql(dbtool.JSONB{"testing": "test123"})
		})

		g.It("Delete the sample post created in DB", func() {
			var post dbtool.Post
			var ID int64
			ID = 1
			dbErr := pool.Delete(&post, map[string]interface{}{"id": ID})
			queryErr := pool.Query(&post, map[string]interface{}{"id": ID})
			g.Assert(dbErr).IsNil()
			g.Assert(queryErr).Eql(errors.New("record not found"))
		})
	})

	tests.RunAuthTests(&toolBox)

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

	// Drop the tables
	pool.Exec("DROP TABLE users")
	pool.Exec("DROP TABLE posts")
	// Need to create users and posts && copy the data from the temp table

	// Note: should separate the test db and production db
}
