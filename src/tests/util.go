package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

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

func loginAndCreatePost(c *Container, p *db.Post, u *userInfo) (*models.Post, []*http.Cookie, error) {
	tx, err := c.DB.BeginTx(c.Context, nil)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Commit()

	user := db.User{Email: u.email, Password: u.pwd, TokenExpiresIn: 0}
	userRecord := db.BindDataToUserModel(&user)
	if err := userRecord.Insert(c.Context, tx, boil.Infer()); err != nil {
		return nil, nil, err
	}
	p.Author = strings.Split(userRecord.Email.String, "@")[0]

	postRecord, err := createPost(c.Context, tx, p)
	if err != nil {
		return nil, nil, err
	}

	loginResult := login(c, u.email, u.pwd)
	return postRecord, loginResult.Result().Cookies(), nil
}

func createPost(ctx context.Context, pool boil.ContextExecutor, p *db.Post) (*models.Post, error) {
	postRecord := db.BindDataToPostModel(p)
	if err := postRecord.Insert(ctx, pool, boil.Infer()); err != nil {
		return nil, err
	}
	return postRecord, nil
}

func createPostsWithRandomTags(c *Container) error {
	tx, err := c.DB.BeginTx(c.Context, nil)
	if err != nil {
		return err
	}
	defer tx.Commit()

	tags := []string{"hello", "nice", "space", "programming", "golang"}
	authors := []string{"Denver", "Mike", "Jessie", "Joe", "Anna"}
	documents := []string{"Something", "Someone", "Somewhere", "Somebody", "Whenever"}
	counter := 0

	p := &db.Post{
		Author: authors[counter],
		Doc:    documents[counter],
		Tags:   tags,
	}

	createPost(c.Context, tx, p)

	for counter < 5 {
		counter += 1

		p := &db.Post{
			Author: authors[createRandomIndex(len(authors)-1)],
			Doc:    documents[createRandomIndex(len(documents)-1)],
			Tags:   tags[:createRandomIndex(len(tags)-1)],
		}
		createPost(c.Context, tx, p)
	}
	return nil
}

func createRandomIndex(max int) int {
	rand.Seed(time.Now().Local().UnixNano())
	min := 0
	limit := max
	return rand.Intn(limit-min+1) + min
}

func verifyCreatedPost(c *Container, result map[string]interface{}, ogPost *db.Post) {
	author, _ := result["author"].(string)
	comments, _ := result["comments"].(string)
	document, _ := result["doc"].(string)
	tags, _ := result["tags"].(string)
	likes, _ := result["likes"].(string)
	likes_int := 0
	if likes != "" {
		likes_int, _ = strconv.Atoi(likes)
	}

	c.Goblin.Assert(author).Eql(strings.Title(ogPost.Author))
	c.Goblin.Assert(comments).Eql(ogPost.Comments)
	c.Goblin.Assert(document).Eql(ogPost.Doc)
	c.Goblin.Assert(likes_int).Eql(ogPost.Likes)
	c.Goblin.Assert(tags).Eql(convertTagsToStr(ogPost.Tags))
}

func extractResult(r *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.Unmarshal(r.Body.Bytes(), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func convertTagsToStr(tags []string) string {

	if len(tags) == 0 {
		return ""
	}

	temp := ""
	for _, s := range tags {
		temp += s + ","
	}

	if len(temp) == 0 {
		return ""
	}

	return temp[:len(temp)-1]
}
