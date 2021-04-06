package tests

import (
	"fmt"
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/db"
)

// RunUsersTests executes all tests for /users
func RunUsersTests(c *Container) {
	c.Goblin.Describe("API /users", func() {
		testCreatUser(c)
		testGetUserWithID(c)
		testUpdateUser(c)
		testDeleteUser(c)
	})
}

func testGetUserWithID(c *Container) {
	c.Goblin.It("/:id GET should return user with given id", func() {
		testUser := createTestUser(c, "test-get-user@test.com", "test-pwd")
		url := fmt.Sprintf("/users/%d", testUser.ID)

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "GET",
			path:    url,
			reqBody: nil,
			cookie:  nil,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()
		verifyCreatedUser(c, response, testUser)
	})

	testGetUserWithInvalidIDType(c)
	testGetUserWithInvalidID(c)
}

func testCreatUser(c *Container) {
	c.Goblin.It("POST should create a new user in database", func() {
		values := Data{
			"email":    "test@test.com",
			"password": "test-password",
		}

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "POST",
			path:    "/users",
			reqBody: &values,
			cookie:  nil,
		})

		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		response, err := extractResult(result)
		c.Goblin.Assert(err).IsNil()
		user, err := db.GetUserByEmail(c.Context, c.DB, "test@test.com")
		c.Goblin.Assert(err).IsNil()
		verifyCreatedUser(c, response, user)
	})

	testCreateUserWithNoPwd(c)

	testCreateUserWithNoEmail(c)

	testCreateUserWithInvalidEmail(c)

	testCreateUserWithNoCredential(c)

	testCreateUserWithInvalidDataType(c)
}

func testUpdateUser(c *Container) {
	c.Goblin.It("PUT should update a user in database", func() {
		testUser := createTestUser(c, "test-update-user@test.com", "test-pwd")
		loginResult := login(c, "test-update-user@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{
			"id":    testUser.ID,
			"email": "something@test.com",
		}

		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "PUT",
			path:    "/users",
			reqBody: &values,
			cookie:  cookies,
		})

		c.Goblin.Assert(result.Code).Eql(http.StatusOK)
		updatedUser := getUserFromDBByID(c, testUser.ID)
		c.Goblin.Assert(values["id"]).Eql(updatedUser.ID)
		c.Goblin.Assert(values["email"]).Eql(updatedUser.Email.String)
	})

	testUpdateUserWithInvalidID(c)

	testUpdateUserWithInvalidEmail(c)

	testUpdateUserWithInvalidDataType(c)

	testUpdateUserWithInvalidCookie(c)

	testUpdateUserWithNoCookie(c)
}

func testDeleteUser(c *Container) {
	testDeleteUserWithInvalidID(c)

	testDeleteUserWithInvalidCookie(c)

	testDeleteUserWithNoCookie(c)

	c.Goblin.It("/:id DELETE should delete a user by its ID", func() {
		testUser := createTestUser(c, "test-delete-user4@test.com", "test-pwd")
		loginResult := login(c, "test-delete-user4@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		reqURL := fmt.Sprintf("/users/%d", testUser.ID)
		result := MakeRequest(&reqData{
			handler: c.Router,
			method:  "DELETE",
			path:    reqURL,
			reqBody: nil,
			cookie:  cookies,
		})
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)
	})
}
