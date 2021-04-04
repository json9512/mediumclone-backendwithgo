package main

import (
	"context"
	"os"
	"testing"

	"github.com/franela/goblin"

	"github.com/json9512/mediumclone-backendwithgo/src/config"
	"github.com/json9512/mediumclone-backendwithgo/src/db"
	"github.com/json9512/mediumclone-backendwithgo/src/tests"
)

func createTestContainer(t *testing.T) *tests.Container {
	config.ReadVariablesFromFile(".env")
	envVars := config.LoadEnvVars()
	logger := config.InitLogger()
	container := db.Init(logger)
	container.Migrate("up")

	router := SetupRouter("test", logger, container.DB)
	g := goblin.Goblin(t)

	testContainer := tests.Container{
		Goblin:  g,
		Router:  router,
		DB:      container.DB,
		Context: context.Background(),
		Env:     envVars,
	}
	return &testContainer
}

func Test(t *testing.T) {
	testContainer := createTestContainer(t)
	defer testContainer.DB.Exec("DROP TABLE gorp_migrations;DROP TABLE users;DROP TABLE posts;")

	tests.RunPostsTests(testContainer)
	tests.RunUsersTests(testContainer)
	tests.RunAuthTests(testContainer)

	// Environment setup test
	testContainer.Goblin.Describe("Environment variables test", func() {
		testContainer.Goblin.It("os.Getenv('DB_NAME') should return $DB_NAME", func() {
			env := os.Getenv("DB_NAME")
			testContainer.Goblin.Assert(env).Equal("mediumclone")
		})
	})
}
