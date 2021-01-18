package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func MakeRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func Test(t *testing.T) {
	// Setup router
	router := SetupRouter("test")

	// Setup test_db for local use only
	// db := db.TestDBInit()

	// create goblin
	g := Goblin(t)
	g.Describe("Server setup test", func() {
		// Passing test
		// GET /ping
		testGetPing(g, router)

		// GET /post
		testGetPost(g, router)

		// GET /post/:id
		testGetPostWithID(g, router)

		// GET /post/:id/like
		testGetLikesOfPost(g, router)

		// GET /post?tag=rabbit
		testGetPostWithQuery(g, router)

		// GET /users
		testGetUsers(g, router)

		// GET /users/:id
		testGetUsersWithID(g, router)
	})

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

}

func testGetPing(g *G, router *gin.Engine) {
	g.It("GET /ping should return JSON {message: pong}", func() {
		// Build expected body
		body := gin.H{
			"message": "pong",
		}

		w := MakeRequest(router, "GET", "/ping")

		g.Assert(w.Code).Equal(http.StatusOK)

		var response map[string]string

		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["message"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["message"]).Equal(value)
	})
}

func testGetPost(g *G, router *gin.Engine) {
	g.It("GET /posts should return list of all posts", func() {
		body := gin.H{
			"result": []string{"test", "sample", "post"},
		}

		w := MakeRequest(router, "GET", "/posts")

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func testGetPostWithID(g *G, router *gin.Engine) {
	g.It("GET /posts/:id should return post with given id", func() {
		body := gin.H{
			"result": "5",
		}

		w := MakeRequest(router, "GET", "/posts/5")

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func testGetLikesOfPost(g *G, router *gin.Engine) {
	g.It("GET /posts/:id/like should return like count of post with given id", func() {
		body := gin.H{
			"result": 10,
		}

		w := MakeRequest(router, "GET", "/posts/5/like")

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]int
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func testGetPostWithQuery(g *G, router *gin.Engine) {
	g.It("GET /posts?tag=rabbit should return tags: [rabbit]", func() {
		tag := make(map[string][]string)
		tag["tags"] = []string{"rabbit"}

		// build expected body
		body := gin.H{
			"result": tag,
		}

		w := MakeRequest(router, "GET", "/posts?tags=rabbit")

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]map[string][]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		// grab the values
		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func testGetUsers(g *G, router *gin.Engine) {
	g.It("GET /users should return list of all users", func() {
		body := gin.H{
			"result": []string{"test", "sample", "users"},
		}

		w := MakeRequest(router, "GET", "/users")

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func testGetUsersWithID(g *G, router *gin.Engine) {
	g.It("GET /users/:id should return user with given id", func() {
		body := gin.H{
			"result": "5",
		}

		w := MakeRequest(router, "GET", "/users/5")

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}
