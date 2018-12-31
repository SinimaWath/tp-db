package database

import (
	"errors"

	"github.com/SinimaWath/tp-db/internal/models"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrUserConflict = errors.New("UserC")
	ErrUserNotFound = errors.New("UserN")
)

// Последовательность Nickname Fullname About Email
func scanUser(r *pgx.Row, user *models.User) error {
	return r.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
}

// Последовательность Nickname Fullname About Email
func scanUserRows(r *pgx.Rows, user *models.User) error {
	return r.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
}
