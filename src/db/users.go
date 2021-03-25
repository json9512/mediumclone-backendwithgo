package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

type User struct {
	Email          string
	Password       string
	TokenExpiresIn int64
}

// GetUserByID retrieves a user by its ID
func GetUserByID(ctx context.Context, db *sql.DB, id int64) (*models.User, error) {
	user, err := models.Users(qm.Where("id = ?", id)).One(ctx, db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail retrieves a user by its email
func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (*models.User, error) {
	user, err := models.Users(qm.Where("email = ?", email)).One(ctx, db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// InsertUser creates a new user in db
func InsertUser(ctx context.Context, db *sql.DB, u *User) (*models.User, error) {
	user := bindDataToUserModel(u)
	fmt.Println(user)

	if err := user.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUserByID deletes the user by its ID
func DeleteUserByID(ctx context.Context, db *sql.DB, id int64) (*models.User, error) {
	user, err := GetUserByID(ctx, db, id)
	if err != nil {
		return nil, err
	}
	if _, err := user.Delete(ctx, db); err != nil {
		return nil, err
	}
	return user, err
}

// UpdateUser updates the user, retrieved by its ID, by the provided user struct
func UpdateUser(ctx context.Context, db *sql.DB, id int64, u *User) (*models.User, error) {
	user, err := GetUserByID(ctx, db, id)
	if err != nil {
		return nil, err
	}
	updateUserModel(user, u)

	if _, err := user.Update(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateTokenExpiresIn updates the TokenExpiresIn column of the given user
func UpdateTokenExpiresIn(ctx context.Context, db *sql.DB, u *models.User, tokenExpiresIn int64) (*models.User, error) {
	u.TokenExpiresIn = null.Int64From(tokenExpiresIn)
	if _, err := u.Update(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	return u, nil
}

func updateUserModel(user *models.User, u *User) {
	if u.Email != "" {
		user.Email = null.StringFrom(u.Email)
	}
	if u.Password != "" {
		user.PWD = null.StringFrom(u.Password)
	}
	if u.TokenExpiresIn > -1 {
		user.TokenExpiresIn = null.Int64From(u.TokenExpiresIn)
	}
}

func bindDataToUserModel(u *User) *models.User {
	return &models.User{
		Email:          null.StringFrom(u.Email),
		PWD:            null.StringFrom(u.Password),
		TokenExpiresIn: null.Int64From(u.TokenExpiresIn),
	}
}
