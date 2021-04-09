package tests

import (
	"fmt"
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

func testGetPostsByTagsAndAuthor(c *Container) {
	c.Goblin.It("?tag=hello%20nice&author=denver GET should return posts with tags=[hello, nice] and author=denver", func() {
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts?tags=hello,nice",
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()

		_, countExists := response["total_count"]
		c.Goblin.Assert(countExists).IsTrue()

		posts, postsExist := response["posts"]
		c.Goblin.Assert(postsExist)

		helloValid := checkIfTagExistsInPosts(posts, "hello")
		c.Goblin.Assert(helloValid).IsTrue()

		niceValid := checkIfTagExistsInPosts(posts, "nice")
		c.Goblin.Assert(niceValid).IsTrue()

		authorValid := checkIfAuthorExistsInPosts(posts, "Denver")
		c.Goblin.Assert(authorValid).IsTrue()
	})
}

func testGetPostsByTags(c *Container) {
	c.Goblin.It("?tag=hello,nice GET should return posts with tags=[hello, nice]", func() {
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts?tags=hello,nice",
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()

		_, countExists := response["total_count"]
		c.Goblin.Assert(countExists).IsTrue()

		posts, postsExist := response["posts"]
		c.Goblin.Assert(postsExist)

		helloValid := checkIfTagExistsInPosts(posts, "hello")
		c.Goblin.Assert(helloValid).IsTrue()

		niceValid := checkIfTagExistsInPosts(posts, "nice")
		c.Goblin.Assert(niceValid).IsTrue()
	})
}

func testGetPostsByTag(c *Container) {
	c.Goblin.It("?tag=hello GET should return posts with tags=[hello]", func() {
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    "/posts?tags=hello",
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()

		_, countExists := response["total_count"]
		c.Goblin.Assert(countExists).IsTrue()

		posts, postsExist := response["posts"]
		c.Goblin.Assert(postsExist)

		isValid := checkIfTagExistsInPosts(posts, "hello")
		c.Goblin.Assert(isValid).IsTrue()
	})
}

func testGetPostLikeWithInvalidID(c *Container) {
	c.Goblin.It("/:id/like GET with invalid ID should return error", func() {
		samplePost := &db.Post{Doc: "Test something", Likes: 1213}
		user := &userInfo{"testing-get-post-likes2@test.com", "test", ""}
		post, cookies, _ := loginAndCreatePost(c, samplePost, user)
		c.makeInvalidReq(&errorTestCase{
			nil,
			"GET",
			fmt.Sprintf("/posts/%d/like", post.ID+1),
			"Post not found.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testGetPostWithInvalidID(c *Container) {
	c.Goblin.It("/:id GET with invalid ID should return error", func() {
		samplePost := &db.Post{Doc: "Test something"}
		user := &userInfo{"testing-get-post2@test.com", "test", ""}
		post, _, _ := loginAndCreatePost(c, samplePost, user)
		c.makeInvalidReq(&errorTestCase{
			nil,
			"GET",
			fmt.Sprintf("/posts/%d", post.ID+1),
			"Post not found.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testCreatePostWithInvalidDoc(c *Container) {
	c.Goblin.It("POST with invalid doc should return error", func() {
		_ = createTestUser(c, "test-badID-post@test.com", "test-pwd")
		loginResult := login(c, "test-badID-post@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"doc": 131313}

		c.makeInvalidReq(&errorTestCase{
			values,
			"POST",
			"/posts",
			"Invalid request data.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testCreatePostWithNoDoc(c *Container) {
	c.Goblin.It("POST with no doc should return error", func() {
		_ = createTestUser(c, "test-noID-post@test.com", "test-pwd")
		loginResult := login(c, "test-noID-post@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{"id": "123"}

		c.makeInvalidReq(&errorTestCase{
			values,
			"POST",
			"/posts",
			"ID, Doc required.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testUpdatePostWithInvalidUser(c *Container) {
	c.Goblin.It("PUT with invalid user should return error", func() {
		u := userInfo{"test-update-post-wrong-author@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, _, _ := loginAndCreatePost(c, &sample, &u)

		createTestUser(c, "test-update-post-no-content@test.com", "test-pwd")
		loginResult := login(c, "test-update-post-no-content@test.com", "test-pwd")

		values := Data{"id": post.ID, "doc": "you are not the author"}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"User is not the author of the post.",
			http.StatusBadRequest,
			loginResult.Result().Cookies(),
		})
	})
}

func testUpdatePostWithInvalidID(c *Container) {
	c.Goblin.It("PUT with invalid post ID should return error", func() {
		u := userInfo{"test-update-post-id@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := loginAndCreatePost(c, &sample, &u)

		values := Data{"id": post.ID + 3}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Post not found.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testUpdatePostWithNoID(c *Container) {
	c.Goblin.It("PUT with no post ID should return error", func() {
		u := userInfo{"test-update-nopost-id@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		_, cookies, _ := loginAndCreatePost(c, &sample, &u)

		values := Data{"doc": "yahoo", "tags": "internet of things"}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"ID required.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testUpdatePostWithNoBody(c *Container) {
	c.Goblin.It("PUT with no request body should return error", func() {
		_ = createTestUser(c, "test-update-nobody@test.com", "test-pwd")
		loginResult := login(c, "test-update-nobody@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		c.makeInvalidReq(&errorTestCase{
			nil,
			"PUT",
			"/posts",
			"ID required.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testUpdatePostWithInvalidDataType(c *Container) {
	c.Goblin.It("PUT with invalid data type should return error", func() {
		u := userInfo{"test-update-datatype@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := loginAndCreatePost(c, &sample, &u)

		values := Data{"id": post.ID, "likes": "abs"}
		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/posts",
			"Invalid request data.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testDeletePostWithInvalidID(c *Container) {
	c.Goblin.It("/:id DELETE with invalid ID should return error", func() {
		u := userInfo{"test-delete-post2@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, cookies, _ := loginAndCreatePost(c, &sample, &u)
		url := fmt.Sprintf("/posts/%d", post.ID+1)

		c.makeInvalidReq(&errorTestCase{
			nil,
			"DELETE",
			url,
			"Post not found.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testDeletePostWithInvalidUser(c *Container) {
	c.Goblin.It("/:id DELETE with invalid user should return error", func() {
		u := userInfo{"test-delete-post3@test.com", "test-pwd", ""}
		sample := db.Post{Doc: ""}
		post, _, _ := loginAndCreatePost(c, &sample, &u)

		_ = createTestUser(c, "test-delete-post4@test.com", "test-pwd")
		loginRes := login(c, "test-delete-post4@test.com", "test-pwd")
		url := fmt.Sprintf("/posts/%d", post.ID)

		c.makeInvalidReq(&errorTestCase{
			nil,
			"DELETE",
			url,
			"User is not the author of the post.",
			http.StatusBadRequest,
			loginRes.Result().Cookies(),
		})
	})
}
