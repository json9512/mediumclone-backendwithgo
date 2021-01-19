package tests

import (
	"encoding/json"
	"net/http"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

// GETPing tests /ping endpoint to retrieve pong
func GETPing(g *goblin.G, router *gin.Engine) {
	g.It("GET /ping should return JSON {message: pong}", func() {
		// Build expected body
		body := Data{
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

// GETPosts tests /posts to retrieve all posts
func GETPosts(g *goblin.G, router *gin.Engine) {
	g.It("GET /posts should return list of all posts", func() {
		body := Data{
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

// GETPostWithID tests /posts/:id to retrieve single post with id
func GETPostWithID(g *goblin.G, router *gin.Engine) {
	g.It("GET /posts/:id should return post with given id", func() {
		body := Data{
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

// GETLikesOfPost tests /posts/:id/like
// to retrieve like count of a single post with given id
func GETLikesOfPost(g *goblin.G, router *gin.Engine) {
	g.It("GET /posts/:id/like should return like count of post with given id", func() {
		body := Data{
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

// GETPostWithQuery tests /posts?queryname=XXX
// to retrieve post/posts based on the query
func GETPostWithQuery(g *goblin.G, router *gin.Engine) {
	g.It("GET /posts?tag=rabbit should return tags: [rabbit]", func() {
		tag := make(map[string][]string)
		tag["tags"] = []string{"rabbit"}

		// build expected body
		body := Data{
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

// POSTPostWithID tests /posts to create a new post in database
func POSTPostWithID(g *goblin.G, router *gin.Engine) {
	g.It("POST /posts should create a new post in database", func() {
		values := Data{"post-id": "5"}
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