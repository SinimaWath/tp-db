package database

import (
	"database/sql"
	"errors"

	"github.com/SinimaWath/tp-db/internal/models"
)

var (
	ErrUserConflict = errors.New("UserC")
	ErrUserNotFound = errors.New("UserN")
)

// Последовательность Nickname Fullname About Email
func scanUser(r *sql.Row, user *models.User) error {
	return r.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
}

// Последовательность Nickname Fullname About Email
func scanUserRows(r *sql.Rows, user *models.User) error {
	return r.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
}
