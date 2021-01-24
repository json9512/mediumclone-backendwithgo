package main

import (
	"encoding/json"
	"errors"
	"net/http"
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

	g.Describe("Authentication/Authorization test", func() {
		g.It("POST /login should attempt to login with the test user", func() {
			// Create sample user before login request
			// and check tokens are empty
			sampleUserData := tests.Data{
				"email":    "login@test.com",
				"password": "login-test-password",
			}

			jsonSampleData, _ := json.Marshal(&sampleUserData)

			createSampleResult := tests.MakeRequest(router, "POST", "/users", jsonSampleData)
			g.Assert(createSampleResult.Code).Eql(http.StatusOK)

			// Fetch the created user
			var user users.User
			dbResult := db.Where("email = ?", "login@test.com").Find(&user)
			g.Assert(dbResult.Error).IsNil()
			g.Assert(user.ID).IsNotNil()
			g.Assert(user.Email).Eql("login@test.com")
			g.Assert(user.AccessToken).Eql("")
			g.Assert(user.RefreshToken).Eql("")

			// Successful login should populate tokens
			postBody := tests.Data{
				"email":    user.Email,
				"password": "login-test-password",
			}

			jsonBody, _ := json.Marshal(&postBody)

			result := tests.MakeRequest(router, "POST", "/login", jsonBody)
			g.Assert(result.Code).Eql(http.StatusOK)

			cookies := result.Result().Cookies()
			accessTokenVal := cookies[0].Value
			refreshTokenVal := cookies[1].Value

			g.Assert(cookies).IsNotNil()
			g.Assert(accessTokenVal).Eql("testing-access-token")
			g.Assert(refreshTokenVal).Eql("testing-refresh-token")

			// Delete the test user
			db.Where("email = ?", "login@test.com").Delete(&user)
		})

		g.It("POST /logout should invalidate token for the user", func() {
			// Create new test user
			user := users.User{
				Email:        "logout@test.com",
				Password:     "logout-test-password",
				AccessToken:  "testing-access-token",
				RefreshToken: "testing-refresh-token",
			}

			jsonBody, _ := json.Marshal(&user)
			result := tests.MakeRequest(router, "POST", "/users", jsonBody)
			g.Assert(result.Code).Eql(http.StatusOK)

			var response map[string]interface{}
			err := json.Unmarshal(result.Body.Bytes(), &response)
			g.Assert(err).IsNil()
			g.Assert(response["email"]).Eql("logout@test.com")

			emptyCookie := []*http.Cookie{}
			cookies := result.Result().Cookies()
			g.Assert(cookies).Eql(emptyCookie)

			// Login with the created user
			loginResult := tests.MakeRequest(router, "POST", "/login", jsonBody)
			g.Assert(loginResult.Code).Eql(http.StatusOK)

			cookies = loginResult.Result().Cookies()
			accessTokenVal := cookies[0].Value
			refreshTokenVal := cookies[1].Value

			g.Assert(cookies).IsNotNil()
			g.Assert(accessTokenVal).Eql("testing-access-token")
			g.Assert(refreshTokenVal).Eql("testing-refresh-token")

			// Test logout from here
			postBody := tests.Data{
				"email": user.Email,
			}

			jsonBody, _ = json.Marshal(&postBody)

			logoutResult := tests.MakeRequest(router, "POST", "/logout", jsonBody)
			g.Assert(logoutResult.Code).Eql(http.StatusOK)

			cookies = logoutResult.Result().Cookies()
			accessTokenVal = cookies[0].Value
			refreshTokenVal = cookies[1].Value

			g.Assert(cookies).IsNotNil()
			g.Assert(accessTokenVal).Eql("")
			g.Assert(refreshTokenVal).Eql("")

			// Soft delete
			db.Where("email = ?", "logout@test.com").Delete(&user)
			// Hard delete
			db.Unscoped().Delete(&user)
		})
	})

	// Drop the users table
	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE posts")

	// Note: should separate the test db and production db
}
