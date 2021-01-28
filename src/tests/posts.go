package tests

import (
	"encoding/json"
	"net/http"
)

// GETPosts tests /posts to retrieve all posts
func GETPosts(tb *TestToolbox) {
	tb.G.It("GET /posts should return list of all posts", func() {
		body := Data{
			"result": []string{"test", "sample", "post"},
		}

		result := MakeRequest(tb.R, "GET", "/posts", nil)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["result"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(body["result"]).Eql(value)
	})
}

// GETPostWithID tests /posts/:id to retrieve single post with id
func GETPostWithID(tb *TestToolbox) {
	tb.G.It("GET /posts/:id should return post with given id", func() {
		body := Data{
			"result": "5",
		}

		result := MakeRequest(tb.R, "GET", "/posts/5", nil)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["result"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(body["result"]).Eql(value)
	})
}

// GETLikesOfPost tests /posts/:id/like
// to retrieve like count of a single post with given id
func GETLikesOfPost(tb *TestToolbox) {
	tb.G.It("GET /posts/:id/like should return like count of post with given id", func() {
		body := Data{
			"result": 10,
		}

		result := MakeRequest(tb.R, "GET", "/posts/5/like", nil)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]int
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["result"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(body["result"]).Eql(value)
	})
}

// GETPostWithQuery tests /posts?queryname=XXX
// to retrieve post/posts based on the query
func GETPostWithQuery(tb *TestToolbox) {
	tb.G.It("GET /posts?tag=rabbit should return tags: [rabbit]", func() {
		tag := make(map[string][]string)
		tag["tags"] = []string{"rabbit"}

		// build expected body
		body := Data{
			"result": tag,
		}

		result := MakeRequest(tb.R, "GET", "/posts?tags=rabbit", nil)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]map[string][]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		// grab the values
		value, exists := response["result"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(body["result"]).Eql(value)
	})
}

// POSTPostWithID tests /posts to create a new post in database
func POSTPostWithID(tb *TestToolbox) {
	tb.G.It("POST /posts should create a new post in database", func() {
		values := Data{"post-id": "5"}

		result := MakeRequest(tb.R, "POST", "/posts", &values)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["post-id"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(values["post-id"]).Eql(value)
	})
}

// PUTPost tests /posts to update a post in database
func PUTPost(tb *TestToolbox) {
	tb.G.It("PUT /posts should update a post in database", func() {
		values := Data{"post-id": "5", "doc": "something"}

		result := MakeRequest(tb.R, "PUT", "/posts", &values)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		postID, IDExists := response["post-id"]
		postDoc, docExists := response["doc"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(IDExists).IsTrue()
		tb.G.Assert(values["post-id"]).Eql(postID)
		tb.G.Assert(docExists).IsTrue()
		tb.G.Assert(values["doc"]).Eql(postDoc)
	})
}

// DELPostWithID tests /posts/:id to delete a post in database
func DELPostWithID(tb *TestToolbox) {
	tb.G.It("DELETE /posts/:id should delete a post with the given ID", func() {
		values := Data{"post-id": "5"}

		result := MakeRequest(tb.R, "DELETE", "/posts", &values)

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		postID, IDExists := response["post-id"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(IDExists).IsTrue()
		tb.G.Assert(values["post-id"]).Eql(postID)
	})
}

// RunPostsTests executes all tests for /posts
func RunPostsTests(toolBox *TestToolbox) {
	toolBox.G.Describe("/posts endpoint tests", func() {

		// GET /posts
		GETPosts(toolBox)

		// GET /posts/:id
		GETPostWithID(toolBox)

		// GET /posts/:id/like
		GETLikesOfPost(toolBox)

		// GET /posts?tag=rabbit
		GETPostWithQuery(toolBox)

		// POST /posts with json {post-id: 5}
		POSTPostWithID(toolBox)

		// PUT /posts with json {post-id: 5, doc: something}
		PUTPost(toolBox)

		// DELETE /posts/:id  with json {post-id: 5}
		DELPostWithID(toolBox)
	})
}
