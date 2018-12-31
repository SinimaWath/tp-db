package database

import (
	"database/sql"

	"github.com/SinimaWath/tp-db/internal/models"
)

const (
	pgErrCodeUniqueViolation  = "23505"
	pgErrForeignKeyViolation  = "23503"
	pgErrCodeNotNullViolation = "23502"
)

const clearQuery = `TRUNCATE ONLY post, vote, thread, forum_user, forum, "user"`

func Clear(db *sql.DB) error {
	_, err := db.Exec(clearQuery)
	return err
}

const statusQuery = `SELECT (SELECT COUNT(*) FROM forum), (SELECT COUNT(*) FROM thread), (SELECT COUNT(*) FROM post), (SELECT COUNT(*) FROM "user")`

func Status(db *sql.DB, s *models.Status) error {
	return db.QueryRow(statusQuery).Scan(&s.Forum, &s.Thread, &s.Post, &s.User)
}
