package database

import (
	"database/sql"
	"errors"
	"time"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/go-openapi/strfmt"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrPostConflict = errors.New("PC")
	ErrPostNotFound = errors.New("PN")
)

// id, author, created, edited, message, parent_id, thread_id, forum_slug
func scanPostRows(r *pgx.Rows, post *models.Post) error {
	parent := sql.NullInt64{}
	created := time.Time{}

	err := r.Scan(&post.ID, &post.Author, &created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if err != nil {
		return err
	}

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}

	date, err := strfmt.ParseDateTime(created.Format(strfmt.MarshalFormat))
	if err != nil {
		post.Created = nil
	} else {
		post.Created = &date
	}

	return err
}

func scanPost(r *pgx.Row, post *models.Post) error {
	parent := sql.NullInt64{}
	created := time.Time{}
	err := r.Scan(&post.ID, &post.Author, &created, &post.IsEdited,
		&post.Message, &parent, &post.Thread, &post.Forum)

	if err != nil {
		return err
	}

	if parent.Valid {
		post.Parent = parent.Int64
	} else {
		post.Parent = 0
	}

	date, err := strfmt.ParseDateTime(created.Format(strfmt.MarshalFormat))
	if err != nil {
		post.Created = nil
	} else {
		post.Created = &date
	}

	return err
}
