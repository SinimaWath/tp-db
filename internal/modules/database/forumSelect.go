package database

import (
	"database/sql"

	"github.com/SinimaWath/tp-db/internal/models"
)

const (
	selectForum = `
	SELECT user_nick, slug, title, thread_count, post_count 
	FROM forum
	WHERE slug = $1
	`
)

func SelectForum(db *sql.DB, forum *models.Forum) error {
	err := scanForum(db.QueryRow(
		selectForum,
		forum.Slug,
	), forum)

	if err == sql.ErrNoRows {
		return ErrForumNotFound
	}

	return err
}
