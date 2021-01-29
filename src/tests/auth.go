package tests

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

type testCred struct {
	userEmail   string
	userPwd     string
	testEmail   string
	testPwd     string
	expectedErr string
}

func testLogin(tb *TestToolbox) {
	tb.G.It("POST /login should attempt to login with the test user", func() {
		// Create sample user before login request
		user := createTestUser(tb, "login@test.com", "test-password")
		// Successful login should populate tokens
		postBody := Data{
			"email":    user.Email,
			"password": "test-password",
		}

		result := MakeRequest(tb.R, "POST", "/login", &postBody)
		tb.G.Assert(result.Code).Eql(http.StatusOK)

		cookies := result.Result().Cookies()
		accessTokenVal := cookies[0].Value

		tb.G.Assert(cookies).IsNotNil()
		tb.G.Assert(accessTokenVal).Eql("testing-access-token")
	})

	tb.G.It("POST /login with invalid credential should return error", func() {
		invalidPwd := testCred{
			userEmail:   "test@test.com",
			userPwd:     "test-pwd",
			testEmail:   "test@test.com",
			testPwd:     "test-pwd-invalid",
			expectedErr: "Authentication failed. Wrong password.",
		}

		loginWithInvalidCred(
			tb,
			invalidPwd,
		)

		invalidEmail := testCred{
			userEmail:   "test1@test.com",
			userPwd:     "test-pwd",
			testEmail:   "test13@test.com",
			testPwd:     "test-pwd",
			expectedErr: "Authentication failed. User does not exist.",
		}

		loginWithInvalidCred(
			tb,
			invalidEmail,
		)
		invalidEmailFormat := testCred{
			userEmail:   "test142@test.com",
			userPwd:     "test-pwd",
			testEmail:   "testtest.com",
			testPwd:     "test-pwd",
			expectedErr: "Authentication failed. Invalid data type.",
		}

		loginWithInvalidCred(
			tb,
			invalidEmailFormat,
		)
	})
}

func testLogout(tb *TestToolbox) {
	tb.G.It("POST /logout should invalidate token for the user", func() {
		// Create new test user
		user := createTestUser(tb, "logout@test.com", "test-password")

		// Login with the created user
		loginResult := MakeRequest(tb.R, "POST", "/login", &user)
		tb.G.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value

		tb.G.Assert(cookies).IsNotNil()
		tb.G.Assert(accessTokenVal).Eql("testing-access-token")

		// Test logout from here
		postBody := Data{
			"email": user.Email,
		}

		logoutResult := MakeRequest(tb.R, "POST", "/logout", &postBody)
		tb.G.Assert(logoutResult.Code).Eql(http.StatusOK)

		cookies = logoutResult.Result().Cookies()
		accessTokenVal = cookies[0].Value

		tb.G.Assert(cookies).IsNotNil()
		tb.G.Assert(accessTokenVal).Eql("")

		// Query the db and check if token is removed
		var userFromDB dbtool.User
		err := tb.P.Query(&userFromDB, map[string]interface{}{"email": "logout@test.com"})
		tb.G.Assert(err).IsNil()
		tb.G.Assert(userFromDB.TokenCreatedAt).Eql((*time.Time)(nil))

	})
}

func createTestUser(tb *TestToolbox, email, pwd string) *dbtool.User {
	// Create sample user before login request
	sampleUserData := Data{
		"email":    email,
		"password": pwd,
	}

	createSampleResult := MakeRequest(tb.R, "POST", "/users", &sampleUserData)
	tb.G.Assert(createSampleResult.Code).Eql(http.StatusOK)

	// Fetch the created user
	var user dbtool.User
	err := tb.P.Query(&user, map[string]interface{}{"email": email})
	tb.G.Assert(err).IsNil()
	tb.G.Assert(user.ID).IsNotNil()
	tb.G.Assert(user.Email).Eql(email)

	return &user
}

func loginWithInvalidCred(tb *TestToolbox, testUser testCred) {
	createTestUser(tb, testUser.userEmail, testUser.userPwd)
	postBody := Data{
		"email":    testUser.testEmail,
		"password": testUser.testPwd,
	}

	result := MakeRequest(tb.R, "POST", "/login", &postBody)
	tb.G.Assert(result.Code).Eql(http.StatusBadRequest)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)
	tb.G.Assert(err).IsNil()

	cookies := result.Result().Cookies()
	tb.G.Assert(len(cookies)).Eql(0)
	tb.G.Assert(response["message"]).Eql(testUser.expectedErr)

}

// RunAuthTests runs test cases for /login and /logout
func RunAuthTests(tb *TestToolbox) {
	tb.G.Describe("Authentication/Authorization test", func() {
		testLogin(tb)
		testLogout(tb)
	})
}
