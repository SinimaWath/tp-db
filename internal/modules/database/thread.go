package database

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/SinimaWath/tp-db/internal/models"
)

var (
	ErrThreadNotFoundAuthorOrForum = errors.New("ThreadNAF")
	ErrThreadNotFound              = errors.New("ThreadN")
	ErrThreadConflict              = errors.New("ThreadC")
)

// Последовательность id, slug, user_nick, created, forum_slug, title, message, votes
func scanThread(r *sql.Row, t *models.Thread) error {
	slug := sql.NullString{}
	err := r.Scan(&t.ID, &slug, &t.Author, &t.Created, &t.Forum, &t.Title, &t.Message, &t.Votes)

	if slug.Valid {
		t.Slug = slug.String
	}

	return err
}

func scanThreadRows(r *sql.Rows, t *models.Thread) error {
	slug := sql.NullString{}
	err := r.Scan(&t.ID, &slug, &t.Author, &t.Created, &t.Forum, &t.Title, &t.Message, &t.Votes)

	if slug.Valid {
		t.Slug = slug.String
	}

	return err
}

func isID(slugOrID string) (int, bool) {
	if value, err := strconv.Atoi(slugOrID); err != nil {
		return -1, false
	} else {
		return value, true
	}
}

func slugToNullable(slug string) sql.NullString {
	nullable := sql.NullString{
		String: slug,
		Valid:  true,
	}
	if slug == "" {
		nullable.Valid = false
	}

	return nullable
}

const (
	checkThreadExistAndGetIDBySlug = `
	SELECT id FROM thread WHERE slug = $1
	`

	checkThreadExistAndGetIDForumSlugBySlug = `
	SELECT id, forum_slug FROM thread WHERE slug = $1
	`

	checkThreadExistAndGetForumSlugByID = `
	SELECT forum_slug FROM thread WHERE id = $1
	`

	checkThreadExistByID = `
	SELECT FROM thread WHERE id = $1
	`
)

func ifThreadExistGetID(db *sql.DB, slug string) (int, bool, error) {
	id := -1
	err := db.QueryRow(checkThreadExistAndGetIDBySlug, slug).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return id, false, nil
		}
		return id, false, err
	}
	return id, true, nil
}

func isThreadExist(db *sql.DB, id int) (bool, error) {
	err := db.QueryRow(checkThreadExistByID, id).Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ifThreadExistAndGetFodumSlugByID(db *sql.DB, id int) (string, bool, error) {
	forum := ""
	err := db.QueryRow(checkThreadExistAndGetForumSlugByID, id).Scan(&forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return forum, false, nil
		}
		return forum, false, err
	}
	return forum, true, nil
}

func ifThreadExistAndGetIDForumSlugBySlug(db *sql.DB, slug string) (string, int, bool, error) {
	id := -1
	forum := ""
	err := db.QueryRow(checkThreadExistAndGetIDForumSlugBySlug, slug).Scan(&id, &forum)
	if err != nil {
		if err == sql.ErrNoRows {
			return forum, id, false, nil
		}
		return forum, id, false, err
	}
	return forum, id, true, nil
}
