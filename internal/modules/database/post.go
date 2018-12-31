package database

import (
	"database/sql"
	"errors"

	"github.com/SinimaWath/tp-db/internal/models"
)

var (
	ErrPostConflict = errors.New("PostC")
	ErrPostNotFound = errors.New("PostN")
)

// id, author, created, edited, message, parent_id, thread_id, forum_slug
func scanPostRows(r *sql.Rows, post *models.Post) error {
	parent := sql.NullInt64{}
	err := r.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}

	return err
}

func scanPost(r *sql.Row, post *models.Post) error {
	parent := sql.NullInt64{}
	err := r.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}

	return err
}
