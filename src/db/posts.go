package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"

	"github.com/json9512/mediumclone-backendwithgo/src/models"
)

// Post contains fields required in a post
type Post struct {
	Author   string
	Doc      string
	Comments string
	Tags     []string
	Likes    int
}

// GetPosts returns all posts
func GetPosts(ctx context.Context, db *sql.DB) (*models.PostSlice, error) {
	posts, err := models.Posts().All(ctx, db)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

// GetPostsByTags returns posts that have a tag listed in the tags array
func GetPostsByTags(ctx context.Context, db *sql.DB, tags []string) (*models.PostSlice, error) {
	posts, err := models.Posts(qm.Where("tags @> ?", pq.Array(tags))).All(ctx, db)
	if err != nil {
		return nil, err
	}
	return &posts, err
}

func GetPostsByAuthor(ctx context.Context, db *sql.DB, author string) (*models.PostSlice, error) {
	posts, err := models.Posts(qm.Where("author = ?", author)).All(ctx, db)
	if err != nil {
		return nil, err
	}
	return &posts, err
}

// GetPostsByTagsAndFilterByAuthor returns posts that have a tag listed in tags array
// and are written by the author
func GetPostsByTagsAndFilterByAuthor(ctx context.Context, db *sql.DB, tags []string, author string) (*models.PostSlice, error) {
	posts, err := models.Posts(qm.Where("tags @> ? AND author = ?", pq.Array(tags), author)).All(ctx, db)
	if err != nil {
		return nil, err
	}
	return &posts, err
}

// GetPostByID returns a post by its ID
func GetPostByID(ctx context.Context, db *sql.DB, id int64) (*models.Post, error) {
	post, err := models.Posts(qm.Where("id = ?", id)).One(ctx, db)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetLikesForPost returns the like count for a post by the given ID
func GetLikesForPost(ctx context.Context, db *sql.DB, id int64) (int, error) {
	post, err := GetPostByID(ctx, db, id)
	if err != nil {
		return 0, err
	}
	return post.Likes.Int, err
}

// InsertPost inserts new post into db with given Post struct
func InsertPost(ctx context.Context, db *sql.DB, p *Post) (*models.Post, error) {
	post := BindDataToPostModel(p)
	if err := post.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	return post, nil
}

// DeletePostByID deletes the post by its ID
func DeletePostByID(ctx context.Context, db *sql.DB, id int64) (*models.Post, error) {
	post, err := GetPostByID(ctx, db, id)
	if err != nil {
		return nil, err
	}
	if _, err = post.Delete(ctx, db); err != nil {
		return nil, err
	}
	return post, nil
}

// UpdatePost updates a post with the provided ID and Post struct
func UpdatePost(ctx context.Context, db *sql.DB, id int64, p *Post) (*models.Post, error) {
	post, err := GetPostByID(ctx, db, id)
	if err != nil {
		return nil, err
	}
	updatePostModel(post, p)

	if _, err := post.Update(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	return post, nil
}

func updatePostModel(post *models.Post, p *Post) {
	if p.Author != "" {
		post.Author = null.StringFrom(p.Author)
	}
	if p.Comments != "" {
		post.Comments = null.StringFrom(p.Comments)
	}
	if p.Doc != "" {
		post.Document = null.StringFrom(p.Doc)
	}
	if p.Likes > 0 {
		post.Likes = null.IntFrom(p.Likes)
	}
	if len(p.Tags) > 0 {
		post.Tags = types.StringArray(p.Tags)
	}
}

func BindDataToPostModel(p *Post) *models.Post {
	return &models.Post{
		Author:   null.StringFrom(p.Author),
		Document: null.StringFrom(p.Doc),
		Likes:    null.IntFrom(p.Likes),
		Comments: null.StringFrom(p.Comments),
		Tags:     types.StringArray(p.Tags),
	}
}
