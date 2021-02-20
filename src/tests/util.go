package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/dbtool"
)

// Data is for structuring req.body/res.body in json format
type Data map[string]interface{}

// TestToolbox contains router, db, and goblin
type TestToolbox struct {
	Goblin *goblin.G
	Router *gin.Engine
	DB     *dbtool.DB
}

type reqData struct {
	handler http.Handler
	method  string
	path    string
	reqBody interface{}
	cookie  []*http.Cookie
}

type errorTestCase struct {
	tb      *TestToolbox
	data    interface{}
	method  string
	url     string
	errMsg  string
	errCode int
	cookies []*http.Cookie
}

// MakeRequest returns the response after making a HTTP request
// with provided parameters
func MakeRequest(r *reqData) *httptest.ResponseRecorder {
	var jsonBody []byte
	if r.reqBody != nil {
		jsonBody, _ = json.Marshal(&r.reqBody)
	}

	req, _ := http.NewRequest(r.method, r.path, bytes.NewBuffer(jsonBody))
	if len(r.cookie) > 0 {
		for _, c := range r.cookie {
			req.AddCookie(c)
		}
	}

	resRecorder := httptest.NewRecorder()
	r.handler.ServeHTTP(resRecorder, req)
	return resRecorder
}

func CreateTestUser(tb *TestToolbox, email, pwd string) *dbtool.User {
	// Create user
	user := Data{
		"email":    email,
		"password": pwd,
	}

	createUserRes := MakeRequest(&reqData{
		handler: tb.Router,
		method:  "POST",
		path:    "/users",
		reqBody: &user,
		cookie:  nil,
	})
	tb.Goblin.Assert(createUserRes.Code).Eql(http.StatusOK)

	// Get user from DB
	testUser, err := tb.DB.GetUserByEmail(email)
	tb.Goblin.Assert(err).IsNil()
	return testUser

}

func LoginUser(tb *TestToolbox, email, pwd string) []*http.Cookie {
	cred := map[string]interface{}{"email": email, "password": pwd}

	loginRes := MakeRequest(&reqData{
		handler: tb.Router,
		method:  "POST",
		path:    "/login",
		reqBody: &cred,
		cookie:  nil,
	})

	tb.Goblin.Assert(loginRes.Code).Eql(http.StatusOK)
	return loginRes.Result().Cookies()
}
