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
	G *goblin.G
	R *gin.Engine
	P *dbtool.Pool
}

// MakeRequest returns the response after making a HTTP request
// with provided parameters
func MakeRequest(r http.Handler, method, path string, reqBody interface{}) *httptest.ResponseRecorder {
	var jsonBody []byte
	if reqBody != nil {
		jsonBody, _ = json.Marshal(&reqBody)
	}
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	resRecorder := httptest.NewRecorder()
	r.ServeHTTP(resRecorder, req)
	return resRecorder
}
