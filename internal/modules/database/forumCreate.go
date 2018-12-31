package database

import (
	"database/sql"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/lib/pq"
)

const (
	insertForum = `
	INSERT INTO forum (user_nick, slug, title) 
	VALUES 
	((SELECT u.nickname FROM "user" u WHERE u.nickname = $1),
	$2, $3)
	RETURNING user_nick, slug, title, thread_count, post_count`
)

func CreateForum(db *sql.DB, forum *models.Forum) error {
	err := scanForum(db.QueryRow(
		insertForum,
		forum.User,
		forum.Slug,
		forum.Title,
	), forum)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError != nil {
			switch pqError.Code {
			case pgErrCodeUniqueViolation:
				return ErrForumConflict
			case pgErrCodeNotNullViolation:
				return ErrForumNotFound
			}
		}
		return err
	}

	return nil
}
