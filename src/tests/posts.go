package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// testGetPosts tests /posts to retrieve all posts
func testGetPosts(c *Container) {
	c.Goblin.It("GET should return list of posts", func() {
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts",
			reqBody: nil,
			cookie:  nil,
		})

		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()

		value, exists := response["posts"]
		c.Goblin.Assert(err).IsNil()
		c.Goblin.Assert(exists).IsTrue()
		c.Goblin.Assert(value).IsNotNil()
	})
}

// testGetPost tests /posts/:id to retrieve a post by id
func testGetPost(c *Container) {
	c.Goblin.It("/:id GET should return a post by its id", func() {
		samplePost := &db.Post{Doc: "Test something"}
		user := &userInfo{"testing-get-post@test.com", "test", ""}
		post, _, _ := loginAndCreatePost(c, samplePost, user)
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    fmt.Sprintf("/posts/%d", post.ID),
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()

		verifyCreatedPost(c, response, samplePost)
	})

	testGetPostWithInvalidID(c)

}

// testGetLikeOfPost tests /posts/:id/like
// to retrieve like count of a post by its id
func testGetLikeOfPost(c *Container) {
	c.Goblin.It("/:id/like GET should return like count of a post by its id", func() {
		samplePost := &db.Post{Doc: "Test something", Likes: 123}
		user := &userInfo{"testing-get-post-likes@test.com", "test", ""}
		post, _, _ := loginAndCreatePost(c, samplePost, user)

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

	testGetPostLikeWithInvalidID(c)
}

// testGetPostWithQuery tests /posts?queryname=XXX
// to retrieve post/posts based on the query
func testGetPostWithQuery(c *Container) {
	testGetPostsByTag(c)

	testGetPostsByTags(c)

	testGetPostsByTagsAndAuthor(c)
}

// testCreatePost tests /posts to create a post in database
func testCreatePost(c *Container) {
	c.Goblin.It("POST should create a post in database", func() {
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

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()

		emptyPost := &db.Post{
			Author:   "Test-Create-Post",
			Doc:      "something",
			Likes:    0,
			Comments: "",
			Tags:     []string{},
		}

		verifyCreatedPost(c, response, emptyPost)
	})

	testCreatePostWithInvalidDoc(c)

	testCreatePostWithNoDoc(c)
}

// testUpdatePost tests /posts to update a post in database
func testUpdatePost(c *Container) {
	c.Goblin.It("PUT should update a post in database", func() {
		u := userInfo{"test-update-post@test.com", "test-pwd", ""}
		sample := &db.Post{Doc: "sample-post"}
		post, cookies, _ := loginAndCreatePost(c, sample, &u)

		values := Data{"id": post.ID, "doc": "something-changed"}
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "PUT",
			path:    "/posts",
			reqBody: &values,
			cookie:  cookies,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)
		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()
		verifyCreatedPost(c, response, sample)
	})

	testUpdatePostWithInvalidUser(c)

	testUpdatePostWithInvalidID(c)

	testUpdatePostWithNoID(c)

	testUpdatePostWithNoBody(c)

	testUpdatePostWithInvalidDataType(c)
}

// testDeletePost tests /posts/:id to delete a post in database
func testDeletePost(c *Container) {
	c.Goblin.It("/:id DELETE should delete a post with the given ID", func() {
		u := userInfo{"test-delete-post@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := loginAndCreatePost(c, &sample, &u)

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

	testDeletePostWithInvalidID(c)

	testDeletePostWithInvalidUser(c)
}

// RunPostsTests executes all tests for /posts
func RunPostsTests(c *Container) {
	c.Goblin.Describe("API /posts", func() {
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
