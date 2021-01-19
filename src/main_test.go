package main

import (
	"os"
	"testing"

	. "github.com/franela/goblin"

	"github.com/json9512/mediumclone-backendwithgo/src/tests"
)

func Test(t *testing.T) {
	// Setup router
	router := SetupRouter("test")

	// Setup test_db for local use only
	// db := db.Init()

	// create goblin
	g := Goblin(t)
	g.Describe("/posts endpoint tests", func() {
		// Passing test
		// GET /ping
		tests.GETPing(g, router)

		// GET /posts
		tests.GETPosts(g, router)

		// GET /posts/:id
		tests.GETPostWithID(g, router)

		// GET /posts/:id/like
		tests.GETLikesOfPost(g, router)

		// GET /posts?tag=rabbit
		tests.GETPostWithQuery(g, router)

		// POST /posts with json {post-id: 5}
		tests.POSTPostWithID(g, router)

		// PUT /posts with json {post-id: 5, doc: something}
		tests.PUTSinglePost(g, router)
	})

	g.Describe("/users endpoint test", func() {
		// GET /users
		tests.GETUsers(g, router)

		// GET /users/:id
		tests.GETUsersWithID(g, router)

		// POST /users with json {user-id: 15}
		tests.POSTUserWithID(g, router)

		// PUT /users with json {user-id: 15, email: something@test.com}
		tests.PUTSingleUser(g, router)
	})

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

}
