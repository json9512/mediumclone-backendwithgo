package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
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

		userID, IDExists := response["id"]
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
		userID, exists := response["id"]
		email, emailExists := response["email"]

		tb.G.Assert(err).IsNil()
		tb.G.Assert(exists).IsTrue()
		tb.G.Assert(emailExists).IsTrue()
		tb.G.Assert(userID).IsNotNil()
		tb.G.Assert(email).Eql(values["email"])
	})

	tb.G.It("POST /users with invalid credential should throw error", func() {
		noPassword := Data{
			"email":    "testUser@test.com",
			"password": "",
		}
		createWithInvalidCred(
			tb,
			noPassword,
			"User registration failed. Invalid credential.")

		noEmail := Data{
			"email":    "",
			"password": "testing",
		}
		createWithInvalidCred(
			tb,
			noEmail,
			"User registration failed. Invalid credential.")

		invalidEmail := Data{
			"email":    "",
			"password": "testing",
		}
		createWithInvalidCred(
			tb,
			invalidEmail,
			"User registration failed. Invalid credential.")

		noCred := Data{
			"email":    "",
			"password": "",
		}
		createWithInvalidCred(
			tb,
			noCred,
			"User registration failed. Invalid credential.")

		invalidData := []string{
			"test@test.com",
			"hello",
		}
		createWithInvalidCred(
			tb,
			invalidData,
			"User registration failed. Invalid data type.")
	})
}

func testUpdateUser(tb *TestToolbox) {
	// Note: the entire failure test case can be refactored for all GET POST PUT DELETE
	tb.G.It("PUT /users should update a user in database", func() {
		var user dbtool.User
		qErr := tb.P.Query(&user, map[string]interface{}{"id": 1})
		tb.G.Assert(qErr).IsNil()

		values := Data{
			"id":    user.ID,
			"email": "something@test.com",
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

		userID, IDExists := response["id"]
		userEmail, emailExists := response["email"]

		// Convert type float64 to uint
		userID = uint(userID.(float64))

		tb.G.Assert(err).IsNil()
		tb.G.Assert(IDExists).IsTrue()
		tb.G.Assert(values["id"]).Eql(userID)
		tb.G.Assert(emailExists).IsTrue()
		tb.G.Assert(values["email"]).Eql(userEmail)
	})

	tb.G.It("PUT /users with invalid ID should return error", func() {
		values := Data{
			"id":    2,
			"email": "something@test.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. Invalid ID.",
			http.StatusBadRequest,
		})
	})

	tb.G.It("PUT /users without new data should return error", func() {
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
		})
	})

	tb.G.It("PUT /users with invalid email should return error", func() {
		var user dbtool.User
		qErr := tb.P.Query(&user, map[string]interface{}{"id": 1})
		tb.G.Assert(qErr).IsNil()

		values := Data{
			"id":    user.ID,
			"email": "somethingtest.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. Invalid email",
			http.StatusBadRequest,
		})
	})

	tb.G.It("PUT /users with invalid data type should return error", func() {
		var user dbtool.User
		qErr := tb.P.Query(&user, map[string]interface{}{"id": 1})
		tb.G.Assert(qErr).IsNil()

		values := []interface{}{
			user.ID,
			"somethingtest.com",
		}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"PUT",
			"/users",
			"User update failed. Invalid data type.",
			http.StatusBadRequest,
		})
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

	tb.G.It("DELETE /users/:id with invalid ID should return error", func() {
		reqURL := fmt.Sprintf("/users/%d", -1)

		makeInvalidReq(&errorTestCase{
			tb,
			nil,
			"DELETE",
			reqURL,
			"Invalid ID",
			http.StatusBadRequest,
		})
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

func createWithInvalidCred(tb *TestToolbox, d interface{}, errorMsg string) {
	result := MakeRequest(&reqData{
		handler: tb.R,
		method:  "POST",
		path:    "/users",
		reqBody: &d,
		cookie:  nil,
	})

	tb.G.Assert(result.Code).Eql(http.StatusBadRequest)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)
	tb.G.Assert(err).IsNil()
	tb.G.Assert(response["message"]).Eql(errorMsg)
}

func makeInvalidReq(e *errorTestCase) {

	result := MakeRequest(&reqData{
		handler: e.tb.R,
		method:  e.method,
		path:    e.url,
		reqBody: &e.data,
		cookie:  nil,
	})

	e.tb.G.Assert(result.Code).Eql(e.errCode)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)

	e.tb.G.Assert(err).IsNil()
	e.tb.G.Assert(response["message"]).Eql(e.errMsg)
}
