package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
)

// Data is for structuring req.body/res.body in json format
type Data map[string]interface{}

// MakeRequest returns the response after making a HTTP request
// with provided parameters
func MakeRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	resRecorder := httptest.NewRecorder()
	r.ServeHTTP(resRecorder, req)

	return resRecorder
}
