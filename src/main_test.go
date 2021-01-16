package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gin-gonic/gin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
)

func MakeRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func Test(t *testing.T) {
	// Setup router
	router := SetupRouter("test")

	// create goblin
	g := Goblin(t)
	g.Describe("Server Test", func() {
		// Passing test
		g.It("GET /ping should return JSON {message: pong}", func() {
			// Build expected body
			body := gin.H{
				"message": "pong",
			}

			// Perform GET request with the handler
			w := MakeRequest(router, "GET", "/ping")

			// Assert we encoded correctly
			// and the request gives 200
			g.Assert(w.Code).Equal(http.StatusOK)

			// Convert JSON response to a map
			var response map[string]string

			err := json.Unmarshal([]byte(w.Body.String()), &response)

			// grab the value
			value, exists := response["message"]

			// make some assertions
			g.Assert(err).IsNil()
			g.Assert(exists).IsTrue()
			g.Assert(body["message"]).Equal(value)
		})

		g.It("GET /posts should return list of all posts", func() {
			// build expected body
			body := gin.H{
				"result": []string{"test", "sample", "post"},
			}

			w := MakeRequest(router, "GET", "/posts")

			g.Assert(w.Code).Eql(http.StatusOK)

			var response map[string][]string
			err := json.Unmarshal([]byte(w.Body.String()), &response)

			// grab the values
			value, exists := response["result"]

			g.Assert(err).IsNil()
			g.Assert(exists).IsTrue()
			g.Assert(body["result"]).Eql(value)
		})
	})

	g.Describe("EnvVar Test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := config.LoadConfig("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})

}
