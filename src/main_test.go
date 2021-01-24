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

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

	// g.Describe("Authentication/Authorization test", func() {
	// 	g.It("POST /login should attempt to login with the test user", func() {
	// 		// Create sample user before login request
	// 		// and check tokens are empty
	// 		var user users.User
	// 		dbResult := db.Where("id = ?", "1").Find(&user)

	// 		g.Assert(dbResult.Error).IsNil()
	// 		g.Assert(user.ID).Eql("1")
	// 		g.Assert(user.Email).Eql("test@test.com")
	// 		g.Assert(user.Password).Eql("test-password")
	// 		g.Assert(user.AccessToken).Eql("")
	// 		g.Assert(user.RefreshToken).Eql("")

	// 		// Successful login should populate tokens
	// 		postBody := tests.Data{
	// 			"email":    user.Email,
	// 			"password": user.Password,
	// 		}

	// 		jsonBody, _ := json.Marshal(&postBody)

	// 		result := tests.MakeRequest(router, "POST", "/login", jsonBody)
	// 		fmt.Println(result.Body.String())
	// 		g.Assert(result.Code).Eql(http.StatusOK)

	// 		var response map[string]string
	// 		err := json.Unmarshal(result.Body.Bytes(), &response)

	// 		accessToken, accessTokenExists := response["access-token"]
	// 		refreshToken, refreshTokenExists := response["refresh-token"]

	// 		g.Assert(err).IsNil()
	// 		g.Assert(accessTokenExists).IsTrue()
	// 		g.Assert(refreshTokenExists).IsTrue()
	// 		g.Assert(accessToken).Eql("testing-access-token")
	// 		g.Assert(refreshToken).Eql("testing-refresh-token")

	// 		// Soft delete
	// 		db.Where("id = ?", "1").Delete(&user)
	// 		// Hard delete
	// 		db.Unscoped().Delete(&user)
	// 	})

	// 	g.It("POST /logout should invalidate token for the user", func() {

	// 		user := users.User{
	// 			Email:        "test@test.com",
	// 			Password:     "test-password",
	// 			AccessToken:  "",
	// 			RefreshToken: "",
	// 		}

	// 		jsonBody, _ := json.Marshal(&user)
	// 		result := tests.MakeRequest(router, "POST", "/login", jsonBody)
	// 		g.Assert(result.Code).Eql(http.StatusOK)

	// 		var response map[string]string
	// 		err := json.Unmarshal(result.Body.Bytes(), &response)
	// 		accessToken, accessTokenExists := response["access-token"]
	// 		refreshToken, refreshTokenExists := response["refresh-token"]

	// 		g.Assert(accessTokenExists).IsTrue()
	// 		g.Assert(refreshTokenExists).IsTrue()
	// 		g.Assert(accessToken).Eql("testing-access-token")
	// 		g.Assert(refreshToken).Eql("testing-refresh-token")

	// 		// Test logout from here
	// 		postBody := tests.Data{
	// 			"email": user.Email,
	// 		}

	// 		jsonBody, _ = json.Marshal(&postBody)

	// 		result = tests.MakeRequest(router, "POST", "/logout", jsonBody)
	// 		g.Assert(result.Code).Eql(http.StatusOK)

	// 		var logoutResponse map[string]string
	// 		err = json.Unmarshal(result.Body.Bytes(), &logoutResponse)

	// 		g.Assert(err).IsNil()
	// 		g.Assert(logoutResponse["access-token"]).Eql("")
	// 		g.Assert(logoutResponse["refresh-token"]).Eql("")

	// 		// Soft delete
	// 		db.Where("email = ?", "test@test.com").Delete(&user)
	// 		// Hard delete
	// 		db.Unscoped().Delete(&user)
	// 	})

	// })

	// Drop the users table
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE posts")

	// Note: should separate the test db and production db
	// for the tests and the server to function properly
	// currently, dropping all the tables after the tests
	// is a hacky way of doing things

}
