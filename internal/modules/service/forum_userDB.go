package service

import (
	"database/sql"
	"errors"

	"github.com/SinimaWath/tp-db/internal/models"
)

var (
	errNotFound       = errors.New("Not found")
	errInternalServer = errors.New("Internal server error")
)

func selectUser(db *sql.DB, user *models.User, nickname string) error {
	querySelect := `SELECT about, email, fullname, nickname FROM "user" WHERE nickname = $1`

	row := db.QueryRow(querySelect, nickname)

	if err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}
