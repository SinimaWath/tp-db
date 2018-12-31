package database

import (
	"database/sql"
	"errors"

	"github.com/SinimaWath/tp-db/internal/models"
)

var (
	ErrForumConflict = errors.New("ForumC")
	ErrForumNotFound = errors.New("ForumN")
)

func scanForum(r *sql.Row, f *models.Forum) error {
	return r.Scan(
		&f.User,
		&f.Slug,
		&f.Title,
		&f.Threads,
		&f.Posts,
	)
}

const (
	checkForumExistQuery = `SELECT FROM forum WHERE slug = $1`
)

func checkForumExist(db *sql.DB, slug string) (bool, error) {
	err := db.QueryRow(checkForumExistQuery, slug).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
