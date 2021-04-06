package tests

import (
	"fmt"
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

func verifyCreatedUser(c *Container, r map[string]interface{}, og *models.User) {
	userID, exists := r["id"]
	userID = int(userID.(float64))
	c.Goblin.Assert(exists).IsTrue()
	c.Goblin.Assert(int(og.ID)).Eql(userID)
	c.Goblin.Assert(og.Email.String).Eql(r["email"])
}

func testGetUserWithInvalidIDType(c *Container) {
	c.Goblin.It("/:id GET with invalid type should return error", func() {
		c.makeInvalidReq(&errorTestCase{
			nil,
			"GET",
			"/users/2@",
			"Invalid ID.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testGetUserWithInvalidID(c *Container) {
	c.Goblin.It("/:id GET with invalid ID should return error", func() {
		c.makeInvalidReq(&errorTestCase{
			nil,
			"GET",
			"/users/-1",
			"Invalid ID.",
			http.StatusBadRequest,
			nil,
		})
	})
}

func testCreateUserWithNoPwd(c *Container) {
	c.Goblin.It("POST with no password should throw error", func() {
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
}

func testCreateUserWithNoEmail(c *Container) {
	c.Goblin.It("POST with no email should throw error", func() {
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
}

func testCreateUserWithInvalidEmail(c *Container) {
	c.Goblin.It("POST with invalid email should throw error", func() {
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
}

func testCreateUserWithNoCredential(c *Container) {
	c.Goblin.It("POST with no credential should throw error", func() {
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
}

func testCreateUserWithInvalidDataType(c *Container) {
	c.Goblin.It("POST with invalid data type should throw error", func() {
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

func testUpdateUserWithInvalidID(c *Container) {
	c.Goblin.It("PUT with invalid ID should return error", func() {
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
			"Invalid request.",
			http.StatusBadRequest,
			cookies,
		})
	})
}

func testUpdateUserWithInvalidEmail(c *Container) {
	c.Goblin.It("PUT with invalid email should return error", func() {
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
}

func testUpdateUserWithInvalidDataType(c *Container) {
	c.Goblin.It("PUT with invalid data type should return error", func() {
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

}

func testUpdateUserWithInvalidCookie(c *Container) {
	c.Goblin.It("PUT with invalid cookie should return error", func() {
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
}

func testUpdateUserWithNoCookie(c *Container) {
	c.Goblin.It("PUT with no cookie should return error", func() {
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

func testDeleteUserWithInvalidID(c *Container) {
	c.Goblin.It("/:id DELETE with invalid ID should return error", func() {
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
}

func testDeleteUserWithInvalidCookie(c *Container) {
	c.Goblin.It("/:id DELETE with invalid cookie should return error", func() {
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
}

func testDeleteUserWithNoCookie(c *Container) {
	c.Goblin.It("/:id DELETE with no cookie should return error", func() {
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
}
