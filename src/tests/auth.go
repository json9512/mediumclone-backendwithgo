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

func testLogin(tb *TestToolbox) {
	tb.Goblin.It("POST /login should attempt to login with the test user", func() {
		// Create sample user before login request
		user := createTestUser(tb, "login@test.com", "test-password")

		// Successful login should populate tokens
		result := login(tb, user.Email, "test-password")
		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		cookies := result.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(accessTokenVal, tb.DB)

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(valid).IsNil()
	})

	tb.Goblin.It("POST /login with invalid password should return error", func() {

		createTestUser(tb, "test-pwd@test.com", "test-pwd")
		result := login(tb, "test-pwd@test.com", "test-pwd-invalid")
		tb.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)
		tb.Goblin.Assert(len(result.Result().Cookies())).Eql(0)

		// Extract error message from result
		body := extractBody(result)
		tb.Goblin.Assert(body["message"]).Eql("Authentication failed. Wrong password.")
	})

	tb.Goblin.It("POST /login with invalid email should return error", func() {
		createTestUser(tb, "test-email@test.com", "test-pwd")
		result := login(tb, "test-email1313@test.com", "test-pwd")
		tb.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)
		tb.Goblin.Assert(len(result.Result().Cookies())).Eql(0)

		// Extract error message from result
		body := extractBody(result)
		tb.Goblin.Assert(body["message"]).Eql("Authentication failed. User does not exist.")
	})

	tb.Goblin.It("POST /login with invalid email format should return error", func() {
		createTestUser(tb, "test-email-2@test.com", "test-pwd")
		result := login(tb, "test-email-2test.com", "test-pwd")
		tb.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)
		tb.Goblin.Assert(len(result.Result().Cookies())).Eql(0)

		// Extract error message from result
		body := extractBody(result)
		tb.Goblin.Assert(body["message"]).Eql("Authentication failed. Invalid data type.")
	})
}

func testLogout(tb *TestToolbox) {
	tb.Goblin.It("POST /logout should invalidate token for the user", func() {
		user := createTestUser(tb, "logout@test.com", "test-password")
		loginResult := login(tb, "logout@test.com", "test-password")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(accessTokenVal, tb.DB)
		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(valid).IsNil()

		// Test logout from here
		logoutResult := logout(tb, user.Email, cookies)

		tb.Goblin.Assert(logoutResult.Code).Eql(http.StatusOK)

		cookies = logoutResult.Result().Cookies()
		accessTokenVal = cookies[0].Value

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(accessTokenVal).Eql("")

		// Query the db and check if token is removed
		userFromDB, err := tb.DB.GetUserByEmail(user.Email)
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFromDB.TokenExpiresIn).Eql(user.TokenExpiresIn)

	})

	tb.Goblin.It("POST /logout with invalid email format should return error", func() {
		createTestUser(tb, "test1422@test.com", "test-pwd")
		loginResult := login(tb, "test1422@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		testAccessToken(tb, loginResult)
		cookies := loginResult.Result().Cookies()

		logoutResult := logout(tb, "testtest.com", cookies)
		tb.Goblin.Assert(logoutResult.Code).Eql(http.StatusBadRequest)

		userFromDB, err := tb.DB.GetUserByEmail("test1422@test.com")
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFromDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		tb.Goblin.Assert(body["message"]).Eql("Logout failed. User does not exist.")

	})

	tb.Goblin.It("POST /logout with invalid email should return error", func() {
		createTestUser(tb, "test131@test.com", "test-pwd")
		loginResult := login(tb, "test131@test.com", "test-pwd")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		testAccessToken(tb, loginResult)
		cookies := loginResult.Result().Cookies()

		logoutResult := logout(tb, "test133@test.com", cookies)
		tb.Goblin.Assert(logoutResult.Code).Eql(http.StatusBadRequest)

		userFromDB, err := tb.DB.GetUserByEmail("test131@test.com")
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFromDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		tb.Goblin.Assert(body["message"]).Eql("Logout failed. User does not exist.")

	})

	tb.Goblin.It("POST /logout with invalid cookie should return error", func() {
		user := createTestUser(tb, "logoutinvalidcookie@test.com", "test-password")
		loginResult := login(tb, "logoutinvalidcookie@test.com", "test-password")
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)
		testAccessToken(tb, loginResult)
		cookies := loginResult.Result().Cookies()

		// mingle the cookie
		cookies[0].Value += "k"
		logoutResult := logout(tb, user.Email, cookies)
		tb.Goblin.Assert(logoutResult.Code).Eql(http.StatusUnauthorized)

		userFromDB, err := tb.DB.GetUserByEmail(user.Email)
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFromDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		tb.Goblin.Assert(body["message"]).Eql("Unauthorized request. Token invalid.")
	})

	tb.Goblin.It("POST /logout with no cookie should return error", func() {
		user := createTestUser(tb, "logoutnocookie@test.com", "test-password")
		// Login with the created user
		loginResult := login(tb, user.Email, user.Password)
		testAccessToken(tb, loginResult)
		logoutResult := logout(tb, user.Email, nil)
		tb.Goblin.Assert(logoutResult.Code).Eql(http.StatusUnauthorized)

		userFromDB, err := tb.DB.GetUserByEmail(user.Email)
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFromDB.TokenExpiresIn).IsNotNil()

		// Extract error message from result
		body := extractBody(logoutResult)
		tb.Goblin.Assert(body["message"]).Eql("Unauthorized request. Token not found.")

	})
}

// RunAuthTests runs test cases for /login and /logout
func RunAuthTests(tb *TestToolbox) {
	tb.Goblin.Describe("Authentication/Authorization test", func() {
		testLogin(tb)
		testLogout(tb)
	})
}
