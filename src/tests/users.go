package tests

import (
	"encoding/json"
	"net/http"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

// GETUsers tests /users to retrieve all users
func GETUsers(g *goblin.G, router *gin.Engine) {
	g.It("GET /users should return list of all users", func() {
		body := Data{
			"result": []string{"test", "sample", "users"},
		}

		w := MakeRequest(router, "GET", "/users", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string][]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

// GETUsers tests /users/:id to retrieve a user by id
func GETUsersWithID(g *goblin.G, router *gin.Engine) {
	g.It("GET /users/:id should return user with given id", func() {
		body := Data{
			"result": "5",
		}

		w := MakeRequest(router, "GET", "/users/5", nil)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["result"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(body["result"]).Eql(value)
	})
}

// POSTUserWithID tests /users to create a new user in database
func POSTUserWithID(g *goblin.G, router *gin.Engine) {
	g.It("POST /users should create a new user in database", func() {
		values := Data{"user-id": "15"}
		jsonValue, _ := json.Marshal(values)

		w := MakeRequest(router, "POST", "/users", jsonValue)

		g.Assert(w.Code).Eql(http.StatusOK)

		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)

		value, exists := response["user-id"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(values["user-id"]).Eql(value)
	})
}
