package main

import (
	"encoding/json"
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

		g.It("users.CreateTestSample should create a sample user in DB", func() {
			var user users.User
			users.CreateTestSample(db)
			result := db.Where("email = ?", "test@test.com").Find(&user)

			g.Assert(result.Error).IsNil()
			g.Assert(user.Email).Eql("test@test.com")
			g.Assert(user.Password).Eql("test-password")
			g.Assert(user.AccessToken).Eql("")
			g.Assert(user.RefreshToken).Eql("")
		})

		g.It("Delete the sample user created in DB", func() {
			var user users.User
			// Soft delete
			db.Where("email = ?", "test@test.com").Delete(&user)

			// Hard delete
			db.Unscoped().Delete(&user)
			result := db.Where("email = ?", "test@test.com").Find(&user)
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
			var user users.User
			users.CreateTestSample(db)
			dbResult := db.Where("email = ?", "test@test.com").Find(&user)

			g.Assert(dbResult.Error).IsNil()
			g.Assert(user.Email).Eql("test@test.com")
			g.Assert(user.Password).Eql("test-password")
			g.Assert(user.AccessToken).Eql("")
			g.Assert(user.RefreshToken).Eql("")

			postBody := tests.Data{
				"email":    user.Email,
				"password": user.Password,
			}
			jsonBody, _ := json.Marshal(&postBody)

			result := tests.MakeRequest(router, "POST", "/login", jsonBody)

			var response map[string]string
			err := json.Unmarshal(result.Body.Bytes(), &response)

			accessToken, accessTokenExists := response["access-token"]
			refreshToken, refreshTokenExists := response["refresh-token"]

			g.Assert(err).IsNil()
			g.Assert(accessTokenExists).IsTrue()
			g.Assert(refreshTokenExists).IsTrue()
			g.Assert(accessToken).Eql("testing-access-token")
			g.Assert(refreshToken).Eql("testing-refresh-token")

			// Soft delete
			db.Where("email = ?", "test@test.com").Delete(&user)
			// Hard delete
			db.Unscoped().Delete(&user)
		})

	})

}
