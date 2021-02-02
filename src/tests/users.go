package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func testGetUserWithID(tb *TestToolbox) {
	tb.G.It("GET /users/:id should return user with given id", func() {
		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "GET",
			path:    "/users/1",
			reqBody: nil,
			cookie:  nil,
		})
		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)

		userID, IDExists := response["user-id"]
		userID = int(userID.(float64))

		tb.G.Assert(err).IsNil()
		tb.G.Assert(IDExists).IsTrue()
		tb.G.Assert(1).Eql(userID)
		tb.G.Assert("test@test.com").Eql(response["email"])
	})

	tb.G.It("GET /users/:id with invalid type should return error", func() {
		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "GET",
			path:    "/users/2@",
			reqBody: nil,
			cookie:  nil,
		})
		tb.G.Assert(result.Code).Eql(http.StatusBadRequest)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		tb.G.Assert(err).IsNil()
		tb.G.Assert(response["message"]).Eql("Invalid ID.")
	})

	tb.G.It("GET /users/:id with invalid ID should return error", func() {
		// Create sample user
		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "GET",
			path:    "/users/-1",
			reqBody: nil,
			cookie:  nil,
		})
		tb.G.Assert(result.Code).Eql(http.StatusBadRequest)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		tb.G.Assert(err).IsNil()
		tb.G.Assert(response["message"]).Eql("User not found.")
	})
}

func testCreatUser(tb *TestToolbox) {
	tb.G.It("POST /users should create a new user in database", func() {
		values := Data{
			"email":    "test@test.com",
			"password": "test-password",
		}

		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "POST",
			path:    "/users",
			reqBody: &values,
			cookie:  nil,
		})

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		userID, exists := response["user-id"]
		email, emailExists := response["email"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(emailExists).IsTrue()
		tb.G.Assert(userID).IsNotNil()
		tb.G.Assert(email).Eql(values["email"])
	})

	tb.G.It("POST /users with invalid credential should throw error", func() {
		createWithInvalidCred(
			tb,
			"testUser@test.com",
			"",
			"User registration failed. Invalid credential.")
		createWithInvalidCred(
			tb,
			"",
			"some-pwd",
			"User registration failed. Invalid credential.")
		createWithInvalidCred(
			tb,
			"invalidemail.com",
			"some-pwd",
			"User registration failed. Invalid credential.")
		createWithInvalidCred(
			tb,
			"",
			"",
			"User registration failed. Invalid credential.")
	})
}

func testUpdateUser(tb *TestToolbox) {
	tb.G.It("PUT /users should update a user in database", func() {
		values := Data{
			"user-id": 1,
			"email":   "something@test.com",
		}

		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "PUT",
			path:    "/users",
			reqBody: &values,
			cookie:  nil,
		})

		tb.G.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)

		userID, IDExists := response["user-id"]
		userEmail, emailExists := response["email"]

		// Convert type float64 to uint
		userID = int(userID.(float64))

		tb.G.Assert(err).IsNil()
		tb.G.Assert(IDExists).IsTrue()
		tb.G.Assert(values["user-id"]).Eql(userID)
		tb.G.Assert(emailExists).IsTrue()
		tb.G.Assert(values["email"]).Eql(userEmail)
	})
}

func testDeleteUser(tb *TestToolbox) {
	tb.G.It("DELETE /users/:id should delete a user with the given ID", func() {
		reqURL := fmt.Sprintf("/users/%d", 1)

		// Perform DELETE request with ID
		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "DELETE",
			path:    reqURL,
			reqBody: nil,
			cookie:  nil,
		})
		tb.G.Assert(result.Code).Eql(http.StatusOK)
	})
}

// RunUsersTests executes all tests for /users
func RunUsersTests(tb *TestToolbox) {
	tb.G.Describe("/users endpoint test", func() {
		testCreatUser(tb)
		testGetUserWithID(tb)
		testUpdateUser(tb)
		testDeleteUser(tb)
	})
}

func createWithInvalidCred(tb *TestToolbox, email, password, errorMsg string) {
	values := Data{
		"email":    email,
		"password": password,
	}

	result := MakeRequest(&reqData{
		handler: tb.R,
		method:  "POST",
		path:    "/users",
		reqBody: &values,
		cookie:  nil,
	})

	tb.G.Assert(result.Code).Eql(http.StatusBadRequest)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)
	tb.G.Assert(err).IsNil()
	tb.G.Assert(response["message"]).Eql(errorMsg)
}
