package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/middlewares"
	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

// Data is for structuring req.body/res.body in json format
type Data map[string]interface{}

// Container contains router, db, and goblin
type Container struct {
	Goblin  *goblin.G
	Router  *gin.Engine
	DB      *sql.DB
	Context context.Context
	Env     *config.EnvVars
}

type reqData struct {
	handler http.Handler
	method  string
	path    string
	reqBody interface{}
	cookie  []*http.Cookie
}

type errorTestCase struct {
	data    interface{}
	method  string
	url     string
	errMsg  string
	errCode int
	cookies []*http.Cookie
}

type userInfo struct {
	email        string
	pwd          string
	access_token string
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

func createTestUser(c *Container, email, pwd string) *models.User {
	// Create user
	user := Data{
		"email":    email,
		"password": pwd,
	}

	createUserRes := MakeRequest(&reqData{
		handler: c.Router,
		method:  "POST",
		path:    "/users",
		reqBody: &user,
		cookie:  nil,
	})
	c.Goblin.Assert(createUserRes.Code).Eql(http.StatusOK)

	createdUser := getUserFromDBByEmail(c, email)
	return createdUser
}

func getUserFromDBByEmail(c *Container, email string) *models.User {
	testUser, err := models.Users(qm.Where("email = ?", email)).One(c.Context, c.DB)
	c.Goblin.Assert(err).IsNil()
	return testUser
}

func getUserFromDBByID(c *Container, id int) *models.User {
	testUser, err := models.Users(qm.Where("id = ?", id)).One(c.Context, c.DB)
	c.Goblin.Assert(err).IsNil()
	return testUser
}

func extractBody(h *httptest.ResponseRecorder) map[string]interface{} {
	var response map[string]interface{}
	_ = json.Unmarshal(h.Body.Bytes(), &response)
	return response
}

func testAccessToken(c *Container, h *httptest.ResponseRecorder) {
	cookies := h.Result().Cookies()
	accessTokenVal := cookies[0].Value
	valid := middlewares.ValidateToken(c.Context, accessTokenVal, c.DB)
	c.Goblin.Assert(cookies).IsNotNil()
	c.Goblin.Assert(valid).IsNil()
}

func login(c *Container, email, password string) *httptest.ResponseRecorder {
	loginBody := Data{
		"email":    email,
		"password": password,
	}

	result := MakeRequest(&reqData{
		handler: c.Router,
		method:  "POST",
		path:    "/login",
		reqBody: &loginBody,
		cookie:  nil,
	})
	return result
}

func logout(c *Container, email string, cookies []*http.Cookie) *httptest.ResponseRecorder {
	data := Data{
		"email": email,
	}

	result := MakeRequest(&reqData{
		handler: c.Router,
		method:  "POST",
		path:    "/logout",
		reqBody: &data,
		cookie:  cookies,
	})
	return result
}

func (c Container) makeInvalidReq(e *errorTestCase) {
	result := MakeRequest(&reqData{
		handler: c.Router,
		method:  e.method,
		path:    e.url,
		reqBody: &e.data,
		cookie:  e.cookies,
	})

	c.Goblin.Assert(result.Code).Eql(e.errCode)

	var response map[string]interface{}
	err := json.Unmarshal(result.Body.Bytes(), &response)

	c.Goblin.Assert(err).IsNil()
	c.Goblin.Assert(response["message"]).Eql(e.errMsg)
}

func createSamplePost(c *Container, p *db.Post, u *userInfo) (*models.Post, []*http.Cookie, error) {
	// use transaction
	tx, err := c.DB.BeginTx(c.Context, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Commit()

	user := db.User{u.email, u.pwd, 0}
	userRecord := db.BindDataToUserModel(&user)
	if err := userRecord.Insert(c.Context, tx, boil.Infer()); err != nil {
		return nil, nil, err
	}
	p.Author = strings.Split(userRecord.Email.String, "@")[0]

	postRecord := db.BindDataToPostModel(p)
	if err := postRecord.Insert(c.Context, tx, boil.Infer()); err != nil {
		return nil, nil, err
	}

	loginResult := login(c, u.email, u.pwd)
	return postRecord, loginResult.Result().Cookies(), nil
}
