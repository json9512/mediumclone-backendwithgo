package tests

import (
	"net/http"

	"github.com/json9512/mediumclone-backendwithgo/src/middlewares"
)

type testCred struct {
	userEmail   string
	userPwd     string
	testEmail   string
	testPwd     string
	expectedErr string
}

func testLogin(c *Container) {
	c.Goblin.It("POST /login should attempt to login with the test user", func() {
		// Create sample user before login request
		user := createTestUser(c, "login@test.com", "test-password")

		// Successful login should populate tokens
		result := login(c, user.Email.String, "test-password")
		c.Goblin.Assert(result.Code).Eql(http.StatusOK)

		cookies := result.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(c.Context, accessTokenVal, c.DB)

		c.Goblin.Assert(cookies).IsNotNil()
		c.Goblin.Assert(valid).IsNil()
	})

	c.Goblin.It("POST /login with invalid password should return error", func() {

		createTestUser(c, "test-pwd@test.com", "test-pwd")
		result := login(c, "test-pwd@test.com", "test-pwd-invalid")
		c.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)
		c.Goblin.Assert(len(result.Result().Cookies())).Eql(0)

		// Extract error message from result
		body := extractBody(result)
		c.Goblin.Assert(body["message"]).Eql("Wrong password.")
	})

	c.Goblin.It("POST /login with invalid email should return error", func() {
		createTestUser(c, "test-email@test.com", "test-pwd")
		result := login(c, "test-email1313@test.com", "test-pwd")
		c.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)
		c.Goblin.Assert(len(result.Result().Cookies())).Eql(0)

		// Extract error message from result
		body := extractBody(result)
		c.Goblin.Assert(body["message"]).Eql("User does not exist.")
	})

	c.Goblin.It("POST /login with invalid email format should return error", func() {
		createTestUser(c, "test-email-2@test.com", "test-pwd")
		result := login(c, "test-email-2test.com", "test-pwd")
		c.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)
		c.Goblin.Assert(len(result.Result().Cookies())).Eql(0)

		// Extract error message from result
		body := extractBody(result)
		c.Goblin.Assert(body["message"]).Eql("Invalid data type.")
	})
}

func testLogout(c *Container) {
	c.Goblin.It("POST /logout should invalidate token for the user", func() {
		user := createTestUser(c, "logout@test.com", "test-password")
		loginResult := login(c, "logout@test.com", "test-password")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(c.Context, accessTokenVal, c.DB)
		c.Goblin.Assert(cookies).IsNotNil()
		c.Goblin.Assert(valid).IsNil()

		// Test logout from here
		logoutResult := logout(c, user.Email.String, cookies)

		c.Goblin.Assert(logoutResult.Code).Eql(http.StatusOK)

		cookies = logoutResult.Result().Cookies()
		accessTokenVal = cookies[0].Value

		c.Goblin.Assert(cookies).IsNotNil()
		c.Goblin.Assert(accessTokenVal).Eql("")

		// Query the db and check if token is removed
		userDB := getUserFromDBByEmail(c, user.Email.String)
		c.Goblin.Assert(userDB.TokenExpiresIn).Eql(user.TokenExpiresIn)

	})

	c.Goblin.It("POST /logout with invalid email format should return error", func() {
		createTestUser(c, "test1422@test.com", "test-pwd")
		loginResult := login(c, "test1422@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		testAccessToken(c, loginResult)
		cookies := loginResult.Result().Cookies()

		logoutResult := logout(c, "testtest.com", cookies)
		c.Goblin.Assert(logoutResult.Code).Eql(http.StatusBadRequest)

		userDB := getUserFromDBByEmail(c, "test1422@test.com")
		c.Goblin.Assert(userDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		c.Goblin.Assert(body["message"]).Eql("User does not exist.")

	})

	c.Goblin.It("POST /logout with invalid email should return error", func() {
		createTestUser(c, "test131@test.com", "test-pwd")
		loginResult := login(c, "test131@test.com", "test-pwd")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		testAccessToken(c, loginResult)
		cookies := loginResult.Result().Cookies()

		logoutResult := logout(c, "test133@test.com", cookies)
		c.Goblin.Assert(logoutResult.Code).Eql(http.StatusBadRequest)

		userDB := getUserFromDBByEmail(c, "test131@test.com")
		c.Goblin.Assert(userDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		c.Goblin.Assert(body["message"]).Eql("User does not exist.")

	})

	c.Goblin.It("POST /logout with invalid cookie should return error", func() {
		user := createTestUser(c, "logoutinvalidcookie@test.com", "test-password")
		loginResult := login(c, "logoutinvalidcookie@test.com", "test-password")
		c.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		testAccessToken(c, loginResult)
		cookies := loginResult.Result().Cookies()

		// mingle the cookie
		cookies[0].Value += "k"
		logoutResult := logout(c, user.Email.String, cookies)
		c.Goblin.Assert(logoutResult.Code).Eql(http.StatusUnauthorized)

		userDB := getUserFromDBByEmail(c, user.Email.String)
		c.Goblin.Assert(userDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		c.Goblin.Assert(body["message"]).Eql("Token invalid.")
	})

	c.Goblin.It("POST /logout with no cookie should return error", func() {
		user := createTestUser(c, "logoutnocookie@test.com", "test-password")
		// Login with the created user
		loginResult := login(c, user.Email.String, user.PWD.String)
		testAccessToken(c, loginResult)
		logoutResult := logout(c, user.Email.String, nil)
		c.Goblin.Assert(logoutResult.Code).Eql(http.StatusUnauthorized)

		userDB := getUserFromDBByEmail(c, user.Email.String)
		c.Goblin.Assert(userDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		c.Goblin.Assert(body["message"]).Eql("Token not found.")

	})
}

// RunAuthTests runs test cases for /login and /logout
func RunAuthTests(c *Container) {
	c.Goblin.Describe("Authentication/Authorization", func() {
		testLogin(c)
		testLogout(c)
	})
}
