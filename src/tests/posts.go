package tests

import (
	"encoding/json"
	"net/http"
)

// testGetPosts tests /posts to retrieve all posts
func testGetPosts(tb *TestToolbox) {
	tb.Goblin.It("GET /posts should return list of all posts", func() {
		body := Data{
			"result": []string{"test", "sample", "post"},
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "GET",
			path:    "/posts",
			reqBody: nil,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["result"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(body["result"]).Eql(value)
	})
}

// testGetPost tests /posts/:id to retrieve single post with id
func testGetPost(tb *TestToolbox) {
	tb.Goblin.It("GET /posts/:id should return post with given id", func() {
		body := Data{
			"result": "5",
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "GET",
			path:    "/posts/5",
			reqBody: nil,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["result"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(body["result"]).Eql(value)
	})
}

// testGetLikeOfPost tests /posts/:id/like
// to retrieve like count of a single post with given id
func testGetLikeOfPost(tb *TestToolbox) {
	tb.Goblin.It("GET /posts/:id/like should return like count of post with given id", func() {
		body := Data{
			"result": 10,
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "GET",
			path:    "/posts/5/like",
			reqBody: nil,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]int
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["result"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(body["result"]).Eql(value)
	})
}

// testGetPostWithQuery tests /posts?queryname=XXX
// to retrieve post/posts based on the query
func testGetPostWithQuery(tb *TestToolbox) {
	tb.Goblin.It("GET /posts?tag=rabbit should return tags: [rabbit]", func() {
		tag := map[string]interface{}{
			"tag": []interface{}{"rabbit"},
		}

		// build expected body
		body := Data{
			"result": tag,
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "GET",
			path:    "/posts?tag=rabbit",
			reqBody: nil,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)

		// grab the values
		value, exists := response["result"]
		bodyTag := body["result"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(bodyTag).Eql(value)
	})
}

// testCreatePost tests /posts to create a new post in database
func testCreatePost(tb *TestToolbox) {
	tb.Goblin.It("POST /posts should create a new post in database", func() {
		values := Data{"id": "5"}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/posts",
			reqBody: &values,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		value, exists := response["id"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(values["id"]).Eql(value)
	})
}

// testUpdatePost tests /posts to update a post in database
func testUpdatePost(tb *TestToolbox) {
	tb.Goblin.It("PUT /posts should update a post in database", func() {
		values := Data{"id": "5", "doc": "something"}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "PUT",
			path:    "/posts",
			reqBody: &values,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		postID, IDExists := response["id"]
		postDoc, docExists := response["doc"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(IDExists).IsTrue()
		tb.Goblin.Assert(values["id"]).Eql(postID)
		tb.Goblin.Assert(docExists).IsTrue()
		tb.Goblin.Assert(values["doc"]).Eql(postDoc)
	})
}

// testDeletePost tests /posts/:id to delete a post in database
func testDeletePost(tb *TestToolbox) {
	tb.Goblin.It("DELETE /posts/:id should delete a post with the given ID", func() {
		values := Data{"id": "5"}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "DELETE",
			path:    "/posts",
			reqBody: &values,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(result.Body.String()), &response)

		postID, IDExists := response["id"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(IDExists).IsTrue()
		tb.Goblin.Assert(values["id"]).Eql(postID)
	})
}

// RunPostsTests executes all tests for /posts
func RunPostsTests(toolBox *TestToolbox) {
	toolBox.Goblin.Describe("/posts endpoint tests", func() {

		// GET /posts
		testGetPosts(toolBox)

		// GET /posts/:id
		testGetPost(toolBox)

		// GET /posts/:id/like
		testGetLikeOfPost(toolBox)

		// GET /posts?tag=rabbit
		testGetPostWithQuery(toolBox)

		// POST /posts with json {id: 5}
		testCreatePost(toolBox)

		// PUT /posts with json {id: 5, doc: something}
		testUpdatePost(toolBox)

		// DELETE /posts/:id  with json {id: 5}
		testDeletePost(toolBox)
	})
}
