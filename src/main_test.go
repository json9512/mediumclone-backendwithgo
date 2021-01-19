package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

type data map[string]interface{}

func MakeRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

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
		GETPing(g, router)

		// GET /posts
		GETPost(g, router)

		// GET /posts/:id
		GETPostWithID(g, router)

		// GET /posts/:id/like
		GETLikesOfPost(g, router)

		// GET /posts?tag=rabbit
		GETPostWithQuery(g, router)

		// POST /posts with json {post-id: 5}
		POSTpostsWithID(g, router)
	})

	g.Describe("/users endpoint test", func() {
		// GET /users
		GETUsers(g, router)

		// GET /users/:id
		GETUsersWithID(g, router)
	})

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

}

func GETPing(g *G, router *gin.Engine) {
	g.It("GET /ping should return JSON {message: pong}", func() {
		// Build expected body
		body := data{
			"message": "pong",
		}

		w := MakeRequest(router, "GET", "/ping", nil)

		g.Assert(w.Code).Equal(http.StatusOK)

		var response map[string]string

		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["message"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["message"]).Equal(value)
	})
}

func GETPost(g *G, router *gin.Engine) {
	g.It("GET /posts should return list of all posts", func() {
		body := data{
			"result": []string{"test", "sample", "post"},
		}

		w := MakeRequest(router, "GET", "/posts", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func GETPostWithID(g *G, router *gin.Engine) {
	g.It("GET /posts/:id should return post with given id", func() {
		body := data{
			"result": "5",
		}

		w := MakeRequest(router, "GET", "/posts/5", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func GETLikesOfPost(g *G, router *gin.Engine) {
	g.It("GET /posts/:id/like should return like count of post with given id", func() {
		body := data{
			"result": 10,
		}

		w := MakeRequest(router, "GET", "/posts/5/like", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]int
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func GETPostWithQuery(g *G, router *gin.Engine) {
	g.It("GET /posts?tag=rabbit should return tags: [rabbit]", func() {
		tag := make(map[string][]string)
		tag["tags"] = []string{"rabbit"}

		// build expected body
		body := data{
			"result": tag,
		}

		w := MakeRequest(router, "GET", "/posts?tags=rabbit", nil)

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

func GETUsers(g *G, router *gin.Engine) {
	g.It("GET /users should return list of all users", func() {
		body := data{
			"result": []string{"test", "sample", "users"},
		}

		w := MakeRequest(router, "GET", "/users", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func GETUsersWithID(g *G, router *gin.Engine) {
	g.It("GET /users/:id should return user with given id", func() {
		body := data{
			"result": "5",
		}

		w := MakeRequest(router, "GET", "/users/5", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

func POSTpostsWithID(g *G, router *gin.Engine) {
	g.It("POST /posts should create a new post in database", func() {
		values := data{"post-id": "5"}
		jsonValue, _ := json.Marshal(values)

		w := MakeRequest(router, "POST", "/posts", jsonValue)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["post-id"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(values["post-id"]).Eql(value)
	})
}
