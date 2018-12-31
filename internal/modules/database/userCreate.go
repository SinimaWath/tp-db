package database

import (
	"database/sql"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/lib/pq"
)

const (
	createUser = `
	INSERT INTO "user" (nickname, fullname, about, email)
	VALUES ($1, $2, $3, $4)
	RETURNING nickname, fullname, about, email
	`

	createForumUserQuery = `
	INSERT INTO forum_user (nickname, forum_slug)
	VALUES ($1, $2)
	ON CONFLICT ON CONSTRAINT unique_forum_user DO NOTHING
	`
)

func CreateUser(db *sql.DB, user *models.User) error {
	err := scanUser(db.QueryRow(
		createUser,
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	), user)

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok && pqError != nil {
			switch pqError.Code {
			case pgErrCodeUniqueViolation:
				return ErrUserConflict
			}
		}
		return err
	}
	return nil
}

func createForumUserTx(tx *sql.Tx, author, forum string) error {
	_, err := tx.Exec(createForumUserQuery, author, forum)
	return err
}
