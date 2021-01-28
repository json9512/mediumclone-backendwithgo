package tests

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

func testLogin(g *goblin.G, router *gin.Engine, p *dbtool.Pool) {
	g.It("POST /login should attempt to login with the test user", func() {
		// Create sample user before login request
		// and check tokens are empty
		sampleUserData := Data{
			"email":    "login@test.com",
			"password": "login-test-password",
		}

		jsonSampleData, _ := json.Marshal(&sampleUserData)

		createSampleResult := MakeRequest(router, "POST", "/users", jsonSampleData)
		g.Assert(createSampleResult.Code).Eql(http.StatusOK)

		// Fetch the created user
		var user dbtool.User
		dbResult := p.Where("email = ?", "login@test.com").Find(&user)
		g.Assert(dbResult.Error).IsNil()
		g.Assert(user.ID).IsNotNil()
		g.Assert(user.Email).Eql("login@test.com")

		// Successful login should populate tokens
		postBody := Data{
			"email":    user.Email,
			"password": "login-test-password",
		}

		jsonBody, _ := json.Marshal(&postBody)

		result := MakeRequest(router, "POST", "/login", jsonBody)
		g.Assert(result.Code).Eql(http.StatusOK)

		cookies := result.Result().Cookies()
		accessTokenVal := cookies[0].Value

		g.Assert(cookies).IsNotNil()
		g.Assert(accessTokenVal).Eql("testing-access-token")
	})
}

func testLogout(g *goblin.G, router *gin.Engine, p *dbtool.Pool) {
	g.It("POST /logout should invalidate token for the user", func() {
		// Create new test user
		user := dbtool.User{
			Email:    "logout@test.com",
			Password: "logout-test-password",
		}

		jsonBody, _ := json.Marshal(&user)
		result := MakeRequest(router, "POST", "/users", jsonBody)
		g.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal(result.Body.Bytes(), &response)
		g.Assert(err).IsNil()
		g.Assert(response["email"]).Eql("logout@test.com")

		emptyCookie := []*http.Cookie{}
		cookies := result.Result().Cookies()
		g.Assert(cookies).Eql(emptyCookie)

		// Login with the created user
		loginResult := MakeRequest(router, "POST", "/login", jsonBody)
		g.Assert(loginResult.Code).Eql(http.StatusOK)

		cookies = loginResult.Result().Cookies()
		accessTokenVal := cookies[0].Value

		g.Assert(cookies).IsNotNil()
		g.Assert(accessTokenVal).Eql("testing-access-token")

		// Test logout from here
		postBody := Data{
			"email": user.Email,
		}

		jsonBody, _ = json.Marshal(&postBody)

		logoutResult := MakeRequest(router, "POST", "/logout", jsonBody)
		g.Assert(logoutResult.Code).Eql(http.StatusOK)

		cookies = logoutResult.Result().Cookies()
		accessTokenVal = cookies[0].Value

		g.Assert(cookies).IsNotNil()
		g.Assert(accessTokenVal).Eql("")

		// Query the db and check if token is removed
		var userFromDB dbtool.User
		err = p.Query(&userFromDB, map[string]interface{}{"email": "logout@test.com"})
		g.Assert(err).IsNil()
		g.Assert(userFromDB.TokenCreatedAt).Eql((*time.Time)(nil))

	})
}

// RunAuthTests runs test cases for /login and /logout
func RunAuthTests(g *goblin.G, router *gin.Engine, db *dbtool.Pool) {
	g.Describe("Authentication/Authorization test", func() {
		testLogin(g, router, db)
		testLogout(g, router, db)
	})
}
