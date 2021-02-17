package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
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
		user := createSampleUser(tb, "login@test.com", "test-password")
		// Successful login should populate tokens
		postBody := Data{
			"email":    user.Email,
			"password": "test-password",
		}

		result := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/login",
			reqBody: &postBody,
			cookie:  nil,
		})
		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		cookies := result.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(accessTokenVal, tb.DB)

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(valid).IsNil()
	})

	tb.Goblin.It("POST /login with invalid credential should return error", func() {
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
	tb.Goblin.It("POST /logout should invalidate token for the user", func() {
		// Create new test user
		user := createSampleUser(tb, "logout@test.com", "test-password")

		// Login with the created user
		loginResult := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/login",
			reqBody: &user,
			cookie:  nil,
		})
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(accessTokenVal, tb.DB)

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(valid).IsNil()

		// Test logout from here
		postBody := Data{
			"email": user.Email,
		}

		logoutResult := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/logout",
			reqBody: &postBody,
			cookie:  cookies,
		})

		tb.Goblin.Assert(logoutResult.Code).Eql(http.StatusOK)

		cookies = logoutResult.Result().Cookies()
		accessTokenVal = cookies[0].Value

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(accessTokenVal).Eql("")

		// Query the db and check if token is removed
		var userFromDB dbtool.User
		err := tb.DB.Query(&userFromDB, map[string]interface{}{"email": user.Email})
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFromDB.TokenExpiryAt).Eql(userFromDB.TokenExpiryAt)

	})

	tb.Goblin.It("POST /logout with invalid email format should return error", func() {
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

	tb.Goblin.It("POST /logout with invalid email should return error", func() {
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
	})

	tb.Goblin.It("POST /logout with invalid cookie should return error", func() {
		user := createSampleUser(tb, "logoutinvalidcookie@test.com", "test-password")
		// Login with the created user
		loginResult := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/login",
			reqBody: &user,
			cookie:  nil,
		})
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(accessTokenVal, tb.DB)

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(valid).IsNil()

		// mingle the cookie
		cookies[0].Value += "k"
		values := Data{"email": user.Email}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"POST",
			"/logout",
			"Unauthorized request. Token invalid.",
			http.StatusUnauthorized,
			cookies,
		})
	})

	tb.Goblin.It("POST /logout with no cookie should return error", func() {
		user := createSampleUser(tb, "logoutnocookie@test.com", "test-password")
		// Login with the created user
		loginResult := MakeRequest(&reqData{
			handler: tb.Router,
			method:  "POST",
			path:    "/login",
			reqBody: &user,
			cookie:  nil,
		})
		tb.Goblin.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies := loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value
		valid := middlewares.ValidateToken(accessTokenVal, tb.DB)

		tb.Goblin.Assert(cookies).IsNotNil()
		tb.Goblin.Assert(valid).IsNil()

		values := Data{"email": user.Email}

		makeInvalidReq(&errorTestCase{
			tb,
			values,
			"POST",
			"/logout",
			"Unauthorized request. Token not found.",
			http.StatusUnauthorized,
			nil,
		})
	})
}

func createSampleUser(tb *TestToolbox, email, pwd string) *dbtool.User {
	// Create sample user before login request
	sampleUserData := Data{
		"email":    email,
		"password": pwd,
	}

	createSampleResult := MakeRequest(&reqData{
		handler: tb.Router,
		method:  "POST",
		path:    "/users",
		reqBody: &sampleUserData,
		cookie:  nil,
	})
	tb.Goblin.Assert(createSampleResult.Code).Eql(http.StatusOK)

	var user dbtool.User
	err := tb.DB.Query(&user, map[string]interface{}{"email": email})
	tb.Goblin.Assert(err).IsNil()
	tb.Goblin.Assert(user.ID).IsNotNil()
	tb.Goblin.Assert(user.Email).Eql(email)

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
			handler: tb.Router,
			method:  "POST",
			path:    "/login",
			reqBody: &loginBody,
			cookie:  nil,
		})
		tb.Goblin.Assert(result.Code).Eql(http.StatusOK)

		valid := middlewares.ValidateToken(result.Result().Cookies()[0].Value, tb.DB)
		tb.Goblin.Assert(valid).IsNil()
		cookies = result.Result().Cookies()
	}

	result = MakeRequest(&reqData{
		handler: tb.Router,
		method:  "POST",
		path:    url,
		reqBody: &postBody,
		cookie:  cookies,
	})
	tb.Goblin.Assert(result.Code).Eql(http.StatusBadRequest)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)
	tb.Goblin.Assert(err).IsNil()

	cookies = result.Result().Cookies()

	if url == "/login" {
		tb.Goblin.Assert(len(cookies)).Eql(0)
	} else if url == "/logout" {
		// user in db should have the token
		var userFrmDB dbtool.User
		err := tb.DB.Query(&userFrmDB, map[string]interface{}{"email": testUser.userEmail})
		tb.Goblin.Assert(err).IsNil()
		tb.Goblin.Assert(userFrmDB.TokenExpiryAt).IsNotNil()
	}

	tb.Goblin.Assert(response["message"]).Eql(testUser.expectedErr)
}

// RunAuthTests runs test cases for /login and /logout
func RunAuthTests(tb *TestToolbox) {
	tb.Goblin.Describe("Authentication/Authorization test", func() {
		testLogin(tb)
		testLogout(tb)
	})
}
