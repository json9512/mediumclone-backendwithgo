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
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

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
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

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
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

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
		_ = createTestUser(tb, "test-create-post@test.com", "test-pwd")
		loginResult := login(tb, "test-create-post@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"doc": "something"}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/posts",
			reqBody: &values,
			cookie:  cookies,
		})
		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		id, exists := response["id"]
		author, exists := response["author"]
		likes, exists := response["likes"]
		document, exists := response["doc"]
		tags, exists := response["tags"]
		comments, exists := response["comments"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(values["doc"]).Eql(document)
		tb.Goblin.Assert(id).IsNotNil()
		tb.Goblin.Assert(author).Eql("test-create-post")
		tb.Goblin.Assert(int(likes.(float64))).Eql(0)
		tb.Goblin.Assert(tags).Eql("")
		tb.Goblin.Assert(comments).Eql("")
	})

	tb.Goblin.It("POST /posts with invalid doc should return error", func() {
		_ = createTestUser(tb, "test-badID-post@test.com", "test-pwd")
		loginResult := login(tb, "test-badID-post@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"doc": 131313}

		tb.makeInvalidReq(&errorTestCase{
			values,
			"POST",
			"/posts",
			"Failed to create post. Invalid data",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("POST /posts with no doc should return error", func() {
		_ = createTestUser(tb, "test-noID-post@test.com", "test-pwd")
		loginResult := login(tb, "test-noID-post@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"id": "123"}

		tb.makeInvalidReq(&errorTestCase{
			values,
			"POST",
			"/posts",
			"Failed to create post. Required information not found: ID, Doc",
			http.StatusBadRequest,
			cookies,
		})
	})
}

// testUpdatePost tests /posts to update a post in database
func testUpdatePost(tb *TestToolbox) {
	tb.Goblin.It("PUT /posts should update a post in database", func() {
		_ = createTestUser(tb, "test-update-post@test.com", "test-pwd")
		loginResult := login(tb, "test-update-post@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		postID := createSamplePost(tb, "sample-post", cookies)

		values := Data{"id": postID, "doc": "something-changed"}
		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "PUT",
			path:    "/posts",
			reqBody: &values,
			cookie:  cookies,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)
	})

	tb.Goblin.It("PUT /posts with no new content should return error", func() {
		_ = createTestUser(tb, "test-update-post-no-content@test.com", "test-pwd")
		loginResult := login(tb, "test-update-post-no-content@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		postID := createSamplePost(tb, "sample-post", cookies)

		values := Data{"id": postID}
		tb.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Failed to update post. No new content",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /posts with invalid user should return error", func() {
		// takes too long
		_ = createTestUser(tb, "test-update-post-wrong-author@test.com", "test-pwd")
		loginResult := login(tb, "test-update-post-wrong-author@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		postID := createSamplePost(tb, "sample-posts", cookies)

		_ = createTestUser(tb, "test-update-post-wrong-author2@test.com", "test-pwd")
		loginResult = login(tb, "test-update-post-wrong-author2@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies = loginResult.Result().Cookies()

		values := Data{"id": postID, "doc": "you are not the author"}
		tb.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Failed to update post. User not post author",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /posts with invalid post ID should return error", func() {
		_ = createTestUser(tb, "test-update-post-id@test.com", "test-pwd")
		loginResult := login(tb, "test-update-post-id@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		postID := createSamplePost(tb, "sample-post", cookies)

		values := Data{"id": postID + 3}
		tb.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Failed to update post. Post not found",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /posts with no post ID should return error", func() {
		_ = createTestUser(tb, "test-update-nopost-id@test.com", "test-pwd")
		loginResult := login(tb, "test-update-nopost-id@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		_ = createSamplePost(tb, "sample-post", cookies)

		values := Data{"doc": "yahoo", "tags": "internet of things"}
		tb.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Failed to update post. Required information not found: ID",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /posts with no request body should return error", func() {
		_ = createTestUser(tb, "test-update-nobody@test.com", "test-pwd")
		loginResult := login(tb, "test-update-nobody@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		tb.makeInvalidReq(&errorTestCase{
			nil,
			"PUT",
			"/posts",
			"Failed to update post. Required information not found: ID",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /posts with invalid data type should return error", func() {
		_ = createTestUser(tb, "test-update-datatype@test.com", "test-pwd")
		loginResult := login(tb, "test-update-datatype@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		postID := createSamplePost(tb, "sample-post", cookies)

		values := Data{"id": postID, "likes": "abs"}
		tb.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Failed to update post. Invalid data",
			http.StatusBadRequest,
			cookies,
		})
	})
}

// testDeletePost tests /posts/:id to delete a post in database
func testDeletePost(tb *TestToolbox) {
	tb.Goblin.It("DELETE /posts/:id should delete a post with the given ID", func() {
		sampleUser := createTestUser(tb, "test-delete-post@test.com", "test-pwd")
		loginResult := login(tb, "test-delete-post@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		values := Data{"id": sampleUser.ID}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "DELETE",
			path:    "/posts",
			reqBody: &values,
			cookie:  cookies,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		postID, IDExists := response["id"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(IDExists).IsTrue()
		tb.Goblin.Assert(values["id"]).Eql(uint(postID.(float64)))
	})
}

// RunPostsTests executes all tests for /posts
func RunPostsTests(toolbox *TestToolbox) {
	toolbox.Goblin.Describe("/posts endpoint tests", func() {

		// GET /posts
		testGetPosts(toolbox)

		// GET /posts/:id
		testGetPost(toolbox)

		// GET /posts/:id/like
		testGetLikeOfPost(toolbox)

		// GET /posts?tag=rabbit
		testGetPostWithQuery(toolbox)

		// POST /posts with json {id: 5}
		testCreatePost(toolbox)

		// PUT /posts with json {id: 5, doc: something}
		testUpdatePost(toolbox)

		// DELETE /posts/:id  with json {id: 5}
		testDeletePost(toolbox)
	})
}
