package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		user := createSampleUser(tb, "login@test.com", "test-password")
		// Successful login should populate tokens
		postBody := Data{
			"email":    user.Email,
			"password": "test-password",
		}

		result := MakeRequest(&reqData{
			handler: tb.R,
			method:  "POST",
			path:    "/login",
			reqBody: &postBody,
			cookie:  nil,
		})
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

		authWithInvalidCred(
			"/login",
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

		authWithInvalidCred(
			"/login",
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

		authWithInvalidCred(
			"/login",
			tb,
			invalidEmailFormat,
		)
	})
}

func testLogout(tb *TestToolbox) {
	tb.G.It("POST /logout should invalidate token for the user", func() {
		// Create new test user
		user := createSampleUser(tb, "logout@test.com", "test-password")

		// Login with the created user
		loginResult := MakeRequest(&reqData{
			handler: tb.R,
			method:  "POST",
			path:    "/login",
			reqBody: &user,
			cookie:  nil,
		})
		tb.G.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value

		tb.G.Assert(cookies).IsNotNil()
		tb.G.Assert(accessTokenVal).Eql("testing-access-token")

		// Test logout from here
		postBody := Data{
			"email": user.Email,
		}

		logoutResult := MakeRequest(&reqData{
			handler: tb.R,
			method:  "POST",
			path:    "/logout",
			reqBody: &postBody,
			cookie:  nil,
		})
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

	tb.G.It("POST /logout with invalid cred should return error", func() {
		invalidEmail := testCred{
			userEmail:   "test131@test.com",
			userPwd:     "test-pwd",
			testEmail:   "test133@test.com",
			testPwd:     "",
			expectedErr: "Logout failed. User does not exist.",
		}

		authWithInvalidCred(
			"/logout",
			tb,
			invalidEmail,
		)

		invalidEmailFormat := testCred{
			userEmail:   "test1422@test.com",
			userPwd:     "test-pwd",
			testEmail:   "testtest.com",
			testPwd:     "",
			expectedErr: "Logout failed. User does not exist.",
		}

		authWithInvalidCred(
			"/logout",
			tb,
			invalidEmailFormat,
		)
	})
}

func createSampleUser(tb *TestToolbox, email, pwd string) *dbtool.User {
	// Create sample user before login request
	sampleUserData := Data{
		"email":    email,
		"password": pwd,
	}

	createSampleResult := MakeRequest(&reqData{
		handler: tb.R,
		method:  "POST",
		path:    "/users",
		reqBody: &sampleUserData,
		cookie:  nil,
	})
	tb.G.Assert(createSampleResult.Code).Eql(http.StatusOK)

	var user dbtool.User
	err := tb.P.Query(&user, map[string]interface{}{"email": email})
	tb.G.Assert(err).IsNil()
	tb.G.Assert(user.ID).IsNotNil()
	tb.G.Assert(user.Email).Eql(email)

	return &user
}

// NOTE: mixing /login and /logout test logic can be confusing
// although they share code. Need to separate the func
// considering sustainability and readability.
func authWithInvalidCred(url string, tb *TestToolbox, testUser testCred) {
	createSampleUser(tb, testUser.userEmail, testUser.userPwd)
	postBody := Data{
		"email":    testUser.testEmail,
		"password": testUser.testPwd,
	}

	var result *httptest.ResponseRecorder
	var cookies []*http.Cookie

	if url == "/logout" {
		// login to create a token
		loginBody := Data{
			"email":    testUser.userEmail,
			"password": testUser.userPwd,
		}

		result = MakeRequest(&reqData{
			handler: tb.R,
			method:  "POST",
			path:    "/login",
			reqBody: &loginBody,
			cookie:  nil,
		})
		tb.G.Assert(result.Code).Eql(http.StatusOK)
		tb.G.Assert(result.Result().Cookies()[0].Value).Eql("testing-access-token")
		cookies = result.Result().Cookies()
	}

	result = MakeRequest(&reqData{
		handler: tb.R,
		method:  "POST",
		path:    url,
		reqBody: &postBody,
		cookie:  cookies,
	})
	tb.G.Assert(result.Code).Eql(http.StatusBadRequest)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)
	tb.G.Assert(err).IsNil()

	cookies = result.Result().Cookies()

	if url == "/login" {
		tb.G.Assert(len(cookies)).Eql(0)
	} else if url == "/logout" {
		// user in db should have the token
		var userFrmDB dbtool.User
		err := tb.P.Query(&userFrmDB, map[string]interface{}{"email": testUser.userEmail})
		tb.G.Assert(err).IsNil()
		tb.G.Assert(userFrmDB.TokenCreatedAt).IsNotNil()
	}

	tb.G.Assert(response["message"]).Eql(testUser.expectedErr)
}

// RunAuthTests runs test cases for /login and /logout
func RunAuthTests(tb *TestToolbox) {
	tb.G.Describe("Authentication/Authorization test", func() {
		testLogin(tb)
		testLogout(tb)
	})
}
