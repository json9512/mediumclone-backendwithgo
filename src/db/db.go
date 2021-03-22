package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

type Container struct {
	DB         *sql.DB
	JWT_SECRET string
}

// Init returns a pointer to a DB container object after connecting to psql db
func Init(l *logrus.Logger) *Container {
	config := createConfig()
	var container Container

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		l.Error(err)
	}

	container.DB = db

	return &Container{db, getEnv("JWT_SECRET", "")}
}

// Migrate performs db migration located in db/migration dir
func (c *Container) Migrate(l *logrus.Logger, method string) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "db/migration",
	}

	migrationMethod := migrate.Up
	if method == "down" {
		migrationMethod = migrate.Down
	}

	n, err := migrate.Exec(c.DB, "postgres", migrations, migrationMethod)
	if err != nil {
		l.Error(err)
	}

	l.Infof("Applied %d migrations\n", n)
	return nil
}
