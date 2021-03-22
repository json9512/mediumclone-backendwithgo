package main

import (
	"os"
	"testing"

	"github.com/franela/goblin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/tests"
)

func Test(t *testing.T) {
	config.ReadVariablesFromFile(".env")
	// Setup test_db for local use only
	logger := config.InitLogger()
	container := db.Init(logger)

	router := SetupRouter("test", logger, container.DB)
	g := goblin.Goblin(t)

	toolBox := tests.TestToolbox{
		Goblin: g,
		Router: router,
		DB:     container.DB,
	}
	tests.RunPostsTests(&toolBox)
	tests.RunUsersTests(&toolBox)
	tests.RunAuthTests(&toolBox)

	// Environment setup test
	g.Describe("Environment variables test", func() {
		g.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			g.Assert(env).Equal("mediumclone")
		})
	})
}
