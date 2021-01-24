package tests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GETUsersWithID tests /users/:id to retrieve a user by id
func GETUsersWithID(g *goblin.G, router *gin.Engine) {
	g.It("GET /users/:id should return user with given id", func() {
		body := Data{
			"user-id": 1,
			"email":   "test@test.com",
		}

		result := MakeRequest(router, "GET", "/users/1", nil)
		g.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		userID, IDExists := response["user-id"]
		userID = int(userID.(float64))

		g.Assert(err).IsNil()
		g.Assert(IDExists).IsTrue()
		g.Assert(body["user-id"]).Eql(userID)
		g.Assert(body["email"]).Eql(response["email"])
	})
}

// POSTUser tests /users to create a new user in database
func POSTUser(g *goblin.G, router *gin.Engine) {
	g.It("POST /users should create a new user in database", func() {
		values := Data{
			"email":    "test@test.com",
			"password": "test-password",
		}
		jsonValue, _ := json.Marshal(values)

		result := MakeRequest(router, "POST", "/users", jsonValue)

		g.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.String()), &response)
		userID, exists := response["user-id"]
		email, emailExists := response["email"]

		g.Assert(err).IsNil()
		g.Assert(exists).IsTrue()
		g.Assert(emailExists).IsTrue()
		g.Assert(userID).IsNotNil()
		g.Assert(email).Eql(values["email"])
	})
}

// PUTSingleUser tests /users to update a user in database
func PUTSingleUser(g *goblin.G, router *gin.Engine) {
	g.It("PUT /users should update a user in database", func() {
		values := Data{
			"user-id": 1,
			"email":   "something@test.com",
		}
		jsonValue, _ := json.Marshal(values)

		result := MakeRequest(router, "PUT", "/users", jsonValue)

		g.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		userID, IDExists := response["user-id"]
		userEmail, emailExists := response["email"]

		// Convert type float64 to uint
		userID = int(userID.(float64))

		g.Assert(err).IsNil()
		g.Assert(IDExists).IsTrue()
		g.Assert(values["user-id"]).Eql(userID)
		g.Assert(emailExists).IsTrue()
		g.Assert(values["email"]).Eql(userEmail)
	})
}

// DELUserWithID tests /users/:id to delete a user in database
func DELUserWithID(g *goblin.G, router *gin.Engine) {
	g.It("DELETE /users/:id should delete a user with the given ID", func() {
		reqURL := fmt.Sprintf("/users/%d", 1)

		// Perform DELETE request with extracted ID
		result := MakeRequest(router, "DELETE", reqURL, nil)
		g.Assert(result.Code).Eql(http.StatusOK)

		var response map[string]interface{}
		err := json.Unmarshal([]byte(result.Body.Bytes()), &response)

		delUserID, IDExists := response["user-id"]
		delUserID = int(delUserID.(float64))

		g.Assert(err).IsNil()
		g.Assert(IDExists).IsTrue()
		g.Assert(delUserID).Eql(1)
		g.Assert(response["email"]).Eql("something@test.com")
		g.Assert(response["access_token"]).Eql("")
		g.Assert(response["refresh_token"]).Eql("")

	})
}

// RunUsersTests executes all tests for /users
func RunUsersTests(g *goblin.G, router *gin.Engine, db *gorm.DB) {
	g.Describe("/users endpoint test", func() {
		POSTUser(g, router)
		GETUsersWithID(g, router)
		PUTSingleUser(g, router)
		DELUserWithID(g, router)
	})
}
