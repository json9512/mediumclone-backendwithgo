package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RunUsersTests executes all tests for /users
func RunUsersTests(tb *TestToolbox) {
	tb.Goblin.Describe("/users endpoint test", func() {
		testCreatUser(tb)
		testGetUserWithID(tb)
		testUpdateUser(tb)
		testDeleteUser(tb)
	})
}

func testGetUserWithID(tb *TestToolbox) {
	tb.Goblin.It("GET /users/:id should return user with given id", func() {
		testUser := createTestUser(tb, "test-get-user@test.com", "test-pwd")
		url := fmt.Sprintf("/users/%d", testUser.ID)

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "GET",
			path:    url,
			reqBody: nil,
			cookie:  nil,
		})
		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)

		userID, IDExists := response["id"]
		userID = int(userID.(float64))

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(IDExists).IsTrue()
		tb.Goblin.Assert(int(testUser.ID)).Eql(userID)
		tb.Goblin.Assert("test-get-user@test.com").Eql(response["email"])
	})

	tb.Goblin.It("GET /users/:id with invalid type should return error", func() {
		makeInvalidReq(&errorTestCase{
			tb,
			nil,
			"GET",
			"/users/2@",
			"Invalid ID.",
			http.StatusBadRequest,
			nil,
		})
	})

	tb.Goblin.It("GET /users/:id with invalid ID should return error", func() {
		makeInvalidReq(&errorTestCase{
			tb,
			nil,
			"GET",
			"/users/-1",
			"User not found.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testCreatUser(tb *TestToolbox) {
	tb.Goblin.It("POST /users should create a new user in database", func() {
		values := Data{
			"email":    "test@test.com",
			"password": "test-password",
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/users",
			reqBody: &values,
			cookie:  nil,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		userID, exists := response["id"]
		email, emailExists := response["email"]

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(exists).IsTrue()
		tb.Goblin.Assert(emailExists).IsTrue()
		tb.Goblin.Assert(userID).IsNotNil()
		tb.Goblin.Assert(email).Eql(values["email"])
	})

	tb.Goblin.It("POST /users with no password should throw error", func() {
		noPassword := Data{
			"email":    "testUser@test.com",
			"password": "",
		}
		makeInvalidReq(&errorTestCase{
			tb,
			noPassword,
			"POST",
			"/users",
			"User registration failed. Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	tb.Goblin.It("POST /users with no email should throw error", func() {
		noEmail := Data{
			"email":    "",
			"password": "testing",
		}
		makeInvalidReq(&errorTestCase{
			tb,
			noEmail,
			"POST",
			"/users",
			"User registration failed. Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	tb.Goblin.It("POST /users with invalid email should throw error", func() {
		invalidEmail := Data{
			"email":    "testtest.com",
			"password": "testing",
		}
		makeInvalidReq(&errorTestCase{
			tb,
			invalidEmail,
			"POST",
			"/users",
			"User registration failed. Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	tb.Goblin.It("POST /users with no credential should throw error", func() {
		noCred := Data{
			"email":    "",
			"password": "",
		}
		makeInvalidReq(&errorTestCase{
			tb,
			noCred,
			"POST",
			"/users",
			"User registration failed. Invalid credential.",
			http.StatusBadRequest,
			nil,
		})
	})

	tb.Goblin.It("POST /users with invalid data type should throw error", func() {
		invalidData := []string{
			"test@test.com",
			"hello",
		}
		makeInvalidReq(&errorTestCase{
			tb,
			invalidData,
			"POST",
			"/users",
			"User registration failed. Invalid data type.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testUpdateUser(tb *TestToolbox) {
	tb.Goblin.It("PUT /users should update a user in database", func() {
		testUser := createTestUser(tb, "test-update-user@test.com", "test-pwd")
		loginResult := login(tb, "test-update-user@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		// attempt to update
		values := Data{
			"id":    testUser.ID,
			"email": "something@test.com",
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "PUT",
			path:    "/users",
			reqBody: &values,
			cookie:  cookies,
		})

		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)

		userID, IDExists := response["id"]
		userEmail, emailExists := response["email"]

		// Convert type float64 to uint
		userID = uint(userID.(float64))

		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(IDExists).IsTrue()
		tb.Goblin.Assert(values["id"]).Eql(userID)
		tb.Goblin.Assert(emailExists).IsTrue()
		tb.Goblin.Assert(values["email"]).Eql(userEmail)
	})

	tb.Goblin.It("PUT /users with invalid ID should return error", func() {
		testUser := createTestUser(tb, "test-put-user@test.com", "test-pwd")
		loginResult := login(tb, "test-put-user@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{
			"id":    testUser.ID + 999,
			"email": "something@test.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. Invalid ID.",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /users without new data should return error", func() {
		_ = createTestUser(tb, "test-put-user2@test.com", "test-pwd")
		loginResult := login(tb, "test-put-user2@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{
			"id": 1,
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. No new data",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /users with invalid email should return error", func() {
		testUser := createTestUser(tb, "test-put-user3@test.com", "test-pwd")
		loginResult := login(tb, "test-put-user3@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := Data{
			"id":    testUser.ID,
			"email": "somethingtest.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. Invalid email",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /users with invalid data type should return error", func() {
		testUser := createTestUser(tb, "test-put-user4@test.com", "test-pwd")
		loginResult := login(tb, "test-put-user4@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		values := []interface{}{
			testUser.ID,
			"somethingtest.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. Invalid data type.",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("PUT /users with invalid cookie should return error", func() {
		testUser := createTestUser(tb, "test-put-user5@test.com", "test-pwd")
		loginResult := login(tb, "test-put-user5@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		cookies[0].Value += "k"

		values := Data{
			"id":    testUser.ID,
			"email": "something@test.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"Unauthorized request. Token invalid.",
			http.StatusUnauthorized,
			cookies,
		})
	})

	tb.Goblin.It("PUT /users with no cookie should return error", func() {
		values := Data{
			"id":    1,
			"email": "something@test.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"Unauthorized request. Token not found.",
			http.StatusUnauthorized,
			nil,
		})
	})
}

func testDeleteUser(tb *TestToolbox) {
	tb.Goblin.It("DELETE /users/:id with invalid ID should return error", func() {
		_ = createTestUser(tb, "test-delete-user@test.com", "test-pwd")
		loginResult := login(tb, "test-delete-user@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		reqURL := fmt.Sprintf("/users/%d", -1)

		makeInvalidReq(&errorTestCase{
			tb,
			nil,
			"DELETE",
			reqURL,
			"Invalid ID",
			http.StatusBadRequest,
			cookies,
		})
	})

	tb.Goblin.It("DELETE /users with invalid cookie should return error", func() {
		testUser := createTestUser(tb, "test-delete-user2@test.com", "test-pwd")
		loginResult := login(tb, "test-delete-user2@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()
		cookies[0].Value += "k"

		reqURL := fmt.Sprintf("/users/%d", testUser.ID)

		makeInvalidReq(&errorTestCase{
			tb,
			nil,
			"DELETE",
			reqURL,
			"Unauthorized request. Token invalid.",
			http.StatusUnauthorized,
			cookies,
		})
	})

	tb.Goblin.It("DELETE /users with no cookie should return error", func() {
		reqURL := fmt.Sprintf("/users/%d", 1)

		makeInvalidReq(&errorTestCase{
			tb,
			nil,
			"DELETE",
			reqURL,
			"Unauthorized request. Token not found.",
			http.StatusUnauthorized,
			nil,
		})
	})

	tb.Goblin.It("DELETE /users/:id should delete a user with the given ID", func() {
		testUser := createTestUser(tb, "test-delete-user4@test.com", "test-pwd")
		loginResult := login(tb, "test-delete-user4@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		cookies := loginResult.Result().Cookies()

		reqURL := fmt.Sprintf("/users/%d", testUser.ID)

		// Perform DELETE request with ID
		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "DELETE",
			path:    reqURL,
			reqBody: nil,
			cookie:  cookies,
		})
		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)
	})
}

func makeInvalidReq(e *errorTestCase) {

	result := MakeRequest(&reqData{
		handler: e.tb.Router,
		method:  e.method,
		path:    e.url,
		reqBody: &e.data,
		cookie:  e.cookies,
	})

	e.tb.Goblin.Assert(result.Code).Eql(e.errCode)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)

	e.tb.Goblin.Assert(err).IsNil()
	e.tb.Goblin.Assert(response["message"]).Eql(e.errMsg)
}
