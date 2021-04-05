package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// testGetPosts tests /posts to retrieve all posts
func testGetPosts(c *Container) {
	c.Goblin.It("GET /posts should return list of all posts", func() {
		body := Data{
			"result": []string{"test", "sample", "post"},
		}

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts",
			reqBody: nil,
			cookie:  nil,
		})

		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		value, exists := response["result"]

		c.Goblin.Assert(err).IsNil()
		c.Goblin.Assert(exists).IsTrue()
		c.Goblin.Assert(body["result"]).Eql(value)
	})
}

// testGetPost tests /posts/:id to retrieve single post with id
func testGetPost(c *Container) {
	// Need to make a post before
	c.Goblin.It("GET /posts/:id should return post with given id", func() {
		samplePost := &db.Post{Doc: "Test something"}
		user := &userInfo{"testing-get-post@test.com", "test", ""}
		post, _, _ := createSamplePost(c, samplePost, user)
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    fmt.Sprintf("/posts/%d", post.ID),
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]string
		_ = json.Unmarshal([]byte(result.Body.Bytes()), &response)
		verifyCreatedPost(c, response, samplePost)
	})
}

// testGetLikeOfPost tests /posts/:id/like
// to retrieve like count of a single post with given id
func testGetLikeOfPost(c *Container) {
	c.Goblin.It("GET /posts/:id/like should return like count of post with given id", func() {
		samplePost := &db.Post{Doc: "Test something", Likes: 123}
		user := &userInfo{"testing-get-post-likes@test.com", "test", ""}
		post, _, _ := createSamplePost(c, samplePost, user)

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    fmt.Sprintf("/posts/%d/like", post.ID),
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]int
		_ = json.Unmarshal([]byte(result.Body.Bytes()), &response)
		c.Goblin.Assert(response["likes"]).Eql(samplePost.Likes)
	})
}

// testGetPostWithQuery tests /posts?queryname=XXX
// to retrieve post/posts based on the query
func testGetPostWithQuery(c *Container) {
	c.Goblin.It("GET /posts?tag=hello should return posts with tags=[hello]", func() {

		// Need to create true case
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts?tags=hello",
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		c.Goblin.Assert(err).IsNil()
		_, countExists := response["totalCount"]
		c.Goblin.Assert(countExists).IsTrue()

		posts, postsExist := response["posts"]
		c.Goblin.Assert(postsExist)

		isValid := checkIfTagExistsInPosts(posts, "hello")
		c.Goblin.Assert(isValid).IsTrue()
	})

	c.Goblin.It("GET /posts?tag=hello,nice should return posts with tags=[hello, nice]", func() {

		// Need to create true case
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts?tags=hello,nice",
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		c.Goblin.Assert(err).IsNil()
		_, countExists := response["totalCount"]
		c.Goblin.Assert(countExists).IsTrue()

		posts, postsExist := response["posts"]
		c.Goblin.Assert(postsExist)

		helloValid := checkIfTagExistsInPosts(posts, "hello")
		c.Goblin.Assert(helloValid).IsTrue()

		niceValid := checkIfTagExistsInPosts(posts, "nice")
		c.Goblin.Assert(niceValid).IsTrue()
	})
}

// testCreatePost tests /posts to create a new post in database
func testCreatePost(c *Container) {
	c.Goblin.It("POST /posts should create a new post in database", func() {
		_ = createTestUser(c, "test-create-post@test.com", "test-pwd")
		loginResult := login(c, "test-create-post@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"doc": "something"}

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "POST",
			path:    "/posts",
			reqBody: &values,
			cookie:  cookies,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		id, exists := response["id"]
		author, exists := response["author"]
		likes, exists := response["likes"]
		document, exists := response["doc"]
		//tags, exists := response["tags"]
		comments, exists := response["comments"]

		c.Goblin.Assert(err).IsNil()
		c.Goblin.Assert(exists).IsTrue()
		c.Goblin.Assert(values["doc"]).Eql(document)
		c.Goblin.Assert(id).IsNotNil()
		c.Goblin.Assert(author).Eql("test-create-post")
		c.Goblin.Assert(int(likes.(float64))).Eql(0)
		//tb.Goblin.Assert(tags.([]interface{})).Eql([]interface{}{""})
		c.Goblin.Assert(comments).Eql("")
	})

	c.Goblin.It("POST /posts with invalid doc should return error", func() {
		_ = createTestUser(c, "test-badID-post@test.com", "test-pwd")
		loginResult := login(c, "test-badID-post@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"doc": 131313}

		c.makeInvalidReq(&errorTestCase{
			values,
			"POST",
			"/posts",
			"Invalid request data",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("POST /posts with no doc should return error", func() {
		_ = createTestUser(c, "test-noID-post@test.com", "test-pwd")
		loginResult := login(c, "test-noID-post@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"id": "123"}

		c.makeInvalidReq(&errorTestCase{
			values,
			"POST",
			"/posts",
			"ID, Doc required",
			http.StatusBadRequest,
			cookies,
		})
	})
}

// testUpdatePost tests /posts to update a post in database
func testUpdatePost(c *Container) {
	c.Goblin.It("PUT /posts should update a post in database", func() {
		u := userInfo{"test-update-post@test.com", "test-pwd", ""}
		sample := db.Post{Doc: "sample-post"}
		post, cookies, _ := createSamplePost(c, &sample, &u)

		values := Data{"id": post.ID, "doc": "something-changed"}
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "PUT",
			path:    "/posts",
			reqBody: &values,
			cookie:  cookies,
		})

		c.Goblin.Assert(result.Code).Eql(http.StatusOK)
	})

	c.Goblin.It("PUT /posts with no new content should return error", func() {
		u := userInfo{"test-update-post-no-content@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := createSamplePost(c, &sample, &u)

		values := Data{"id": post.ID}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"No new content",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /posts with invalid user should return error", func() {
		u := userInfo{"test-update-post-wrong-author@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, _, _ := createSamplePost(c, &sample, &u)
		loginResult := login(c, "test-update-post-no-content@test.com", "test-pwd")

		values := Data{"id": post.ID, "doc": "you are not the author"}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"User is not the author of the post",
			http.StatusBadRequest,
			loginResult.Result().Cookies(),
		})
	})

	c.Goblin.It("PUT /posts with invalid post ID should return error", func() {
		u := userInfo{"test-update-post-id@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := createSamplePost(c, &sample, &u)

		values := Data{"id": post.ID + 3}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Post not found",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /posts with no post ID should return error", func() {
		u := userInfo{"test-update-nopost-id@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		_, cookies, _ := createSamplePost(c, &sample, &u)

		values := Data{"doc": "yahoo", "tags": "internet of things"}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"ID required",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /posts with no request body should return error", func() {
		_ = createTestUser(c, "test-update-nobody@test.com", "test-pwd")
		loginResult := login(c, "test-update-nobody@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		c.makeInvalidReq(&errorTestCase{
			nil,
			"PUT",
			"/posts",
			"ID required",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /posts with invalid data type should return error", func() {
		u := userInfo{"test-update-datatype@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := createSamplePost(c, &sample, &u)

		values := Data{"id": post.ID, "likes": "abs"}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Invalid request data",
			http.StatusBadRequest,
			cookies,
		})
	})
}

// testDeletePost tests /posts/:id to delete a post in database
func testDeletePost(c *Container) {
	c.Goblin.It("DELETE /posts/:id should delete a post with the given ID", func() {
		u := userInfo{"test-delete-post@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := createSamplePost(c, &sample, &u)

		// Delete
		url := fmt.Sprintf("/posts/%d", post.ID)

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "DELETE",
			path:    url,
			reqBody: nil,
			cookie:  cookies,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

	})
}

// RunPostsTests executes all tests for /posts
func RunPostsTests(c *Container) {
	c.Goblin.Describe("/posts endpoint tests", func() {
		createPostsWithRandomTags(c)
		// GET /posts
		testGetPosts(c)

		// GET /posts/:id
		testGetPost(c)

		// GET /posts/:id/like
		testGetLikeOfPost(c)

		// GET /posts?tag=rabbit
		testGetPostWithQuery(c)

		// POST /posts with json {id: 5}
		testCreatePost(c)

		// PUT /posts with json {id: 5, doc: something}
		testUpdatePost(c)

		// DELETE /posts/:id  with json {id: 5}
		testDeletePost(c)
	})
}

func checkIfTagExistsInPosts(posts interface{}, tag string) bool {
	for _, p := range posts.([]interface{}) {
		postTags := p.(map[string]interface{})["tags"].([]interface{})

		for _, c := range postTags {
			if c.(string) == tag {
				return true
			}
		}
		return false
	}
	return false
}
