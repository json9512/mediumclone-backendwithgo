package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RunUsersTests executes all tests for /users
func RunUsersTests(c *Container) {
	c.Goblin.Describe("API", func() {
		testCreatUser(c)
		testGetUserWithID(c)
		testUpdateUser(c)
		testDeleteUser(c)
	})
}

func testGetUserWithID(c *Container) {
	c.Goblin.It("GET /users/:id should return user with given id", func() {
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

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)

		userID, IDExists := response["id"]
		userID = int(userID.(float64))

		c.Goblin.Assert(err).IsNil()
		c.Goblin.Assert(IDExists).IsTrue()
		c.Goblin.Assert(int(testUser.ID)).Eql(userID)
		c.Goblin.Assert("test-get-user@test.com").Eql(response["email"])
	})

	c.Goblin.It("GET /users/:id with invalid type should return error", func() {
		c.makeInvalidReq(&errorTestCase{
			nil,
			"GET",
			"/users/2@",
			"Invalid ID.",
			http.StatusBadRequest,
			nil,
		})
	})

	c.Goblin.It("GET /users/:id with invalid ID should return error", func() {
		c.makeInvalidReq(&errorTestCase{
			nil,
			"GET",
			"/users/-1",
			"User not found.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testCreatUser(c *Container) {
	c.Goblin.It("POST /users should create a new user in database", func() {
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

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		userID, exists := response["id"]
		email, emailExists := response["email"]

		c.Goblin.Assert(err).IsNil()
		c.Goblin.Assert(exists).IsTrue()
		c.Goblin.Assert(emailExists).IsTrue()
		c.Goblin.Assert(userID).IsNotNil()
		c.Goblin.Assert(email).Eql(values["email"])
	})

	c.Goblin.It("POST /users with no password should throw error", func() {
		noPassword := Data{
			"email":    "testUser@test.com",
			"password": "",
		}
		c.makeInvalidReq(&errorTestCase{
			noPassword,
			"POST",
			"/users",
			"Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	c.Goblin.It("POST /users with no email should throw error", func() {
		noEmail := Data{
			"email":    "",
			"password": "testing",
		}
		c.makeInvalidReq(&errorTestCase{
			noEmail,
			"POST",
			"/users",
			"Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	c.Goblin.It("POST /users with invalid email should throw error", func() {
		invalidEmail := Data{
			"email":    "testtest.com",
			"password": "testing",
		}
		c.makeInvalidReq(&errorTestCase{
			invalidEmail,
			"POST",
			"/users",
			"Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	c.Goblin.It("POST /users with no credential should throw error", func() {
		noCred := Data{
			"email":    "",
			"password": "",
		}
		c.makeInvalidReq(&errorTestCase{
			noCred,
			"POST",
			"/users",
			"Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	c.Goblin.It("POST /users with invalid data type should throw error", func() {
		invalidData := []string{
			"test@test.com",
			"hello",
		}
		c.makeInvalidReq(&errorTestCase{
			invalidData,
			"POST",
			"/users",
			"Invalid data type.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testUpdateUser(c *Container) {
	c.Goblin.It("PUT /users should update a user in database", func() {
		testUser := createTestUser(c, "test-update-user@test.com", "test-pwd")
		loginResult := login(c, "test-update-user@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		// attempt to update
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

	c.Goblin.It("PUT /users with invalid ID should return error", func() {
		testUser := createTestUser(c, "test-put-user@test.com", "test-pwd")
		loginResult := login(c, "test-put-user@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{
			"id":    testUser.ID + 999,
			"email": "something@test.com",
		}

		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/users",
			"Invalid ID.",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /users without new data should return error", func() {
		_ = createTestUser(c, "test-put-user2@test.com", "test-pwd")
		loginResult := login(c, "test-put-user2@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		fmt.Println("test")
		values := Data{
			"id": 1,
		}

		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/users",
			"No new data",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /users with invalid email should return error", func() {
		testUser := createTestUser(c, "test-put-user3@test.com", "test-pwd")
		loginResult := login(c, "test-put-user3@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{
			"id":    testUser.ID,
			"email": "somethingtest.com",
		}

		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/users",
			"Invalid data.",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /users with invalid data type should return error", func() {
		testUser := createTestUser(c, "test-put-user4@test.com", "test-pwd")
		loginResult := login(c, "test-put-user4@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := []interface{}{
			testUser.ID,
			"somethingtest.com",
		}

		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/users",
			"Invalid data type.",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("PUT /users with invalid cookie should return error", func() {
		testUser := createTestUser(c, "test-put-user5@test.com", "test-pwd")
		loginResult := login(c, "test-put-user5@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		cookies[0].Value += "k"

		values := Data{
			"id":    testUser.ID,
			"email": "something@test.com",
		}

		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/users",
			"Token invalid.",
			http.StatusUnauthorized,
			cookies,
		})
	})

	c.Goblin.It("PUT /users with no cookie should return error", func() {
		values := Data{
			"id":    1,
			"email": "something@test.com",
		}

		c.makeInvalidReq(&errorTestCase{
			values,
			"PUT",
			"/users",
			"Token not found.",
			http.StatusUnauthorized,
			nil,
		})
	})
}

func testDeleteUser(c *Container) {
	c.Goblin.It("DELETE /users/:id with invalid ID should return error", func() {
		_ = createTestUser(c, "test-delete-user@test.com", "test-pwd")
		loginResult := login(c, "test-delete-user@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		reqURL := fmt.Sprintf("/users/%d", -1)

		c.makeInvalidReq(&errorTestCase{
			nil,
			"DELETE",
			reqURL,
			"Invalid ID.",
			http.StatusBadRequest,
			cookies,
		})
	})

	c.Goblin.It("DELETE /users with invalid cookie should return error", func() {
		testUser := createTestUser(c, "test-delete-user2@test.com", "test-pwd")
		loginResult := login(c, "test-delete-user2@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		cookies[0].Value += "k"

		reqURL := fmt.Sprintf("/users/%d", testUser.ID)

		c.makeInvalidReq(&errorTestCase{
			nil,
			"DELETE",
			reqURL,
			"Token invalid.",
			http.StatusUnauthorized,
			cookies,
		})
	})

	c.Goblin.It("DELETE /users with no cookie should return error", func() {
		reqURL := fmt.Sprintf("/users/%d", 1)

		c.makeInvalidReq(&errorTestCase{
			nil,
			"DELETE",
			reqURL,
			"Token not found.",
			http.StatusUnauthorized,
			nil,
		})
	})

	c.Goblin.It("DELETE /users/:id should delete a user with the given ID", func() {
		testUser := createTestUser(c, "test-delete-user4@test.com", "test-pwd")
		loginResult := login(c, "test-delete-user4@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		reqURL := fmt.Sprintf("/users/%d", testUser.ID)

		// Perform DELETE request with ID
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
