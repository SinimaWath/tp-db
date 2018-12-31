package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/go-openapi/strfmt"
	pgx "gopkg.in/jackc/pgx.v2"
)

var (
	ErrThreadNotFoundAuthorOrForum = errors.New("ThreadNAF")
	ErrThreadNotFound              = errors.New("ThreadN")
	ErrThreadConflict              = errors.New("ThreadC")
)

// Последовательность id, slug, user_nick, created, forum_slug, title, message, votes
func scanThread(r *pgx.Row, t *models.Thread) error {
	slug := sql.NullString{}
	created := pgx.NullTime{}
	err := r.Scan(&t.ID, &slug, &t.Author, &created, &t.Forum, &t.Title, &t.Message, &t.Votes)
	if err != nil {
		return err
	}
	if slug.Valid {
		t.Slug = slug.String
	}

	if created.Valid {
		date, err := strfmt.ParseDateTime(created.Time.Format(strfmt.MarshalFormat))
		if err != nil {
			t.Created = nil
		} else {
			t.Created = &date
		}
	} else {
		t.Created = nil
	}

	return err
}

func scanThreadRows(r *pgx.Rows, t *models.Thread) error {
	slug := sql.NullString{}
	created := time.Time{}
	err := r.Scan(&t.ID, &slug, &t.Author, &created, &t.Forum, &t.Title, &t.Message, &t.Votes)
	if err != nil {
		return err
	}
	if slug.Valid {
		t.Slug = slug.String
	}

	date, err := strfmt.ParseDateTime(created.Format(strfmt.MarshalFormat))
	if err != nil {
		t.Created = nil
	} else {
		t.Created = &date
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

func ifThreadExistGetID(db *pgx.ConnPool, slug string) (int, bool, error) {
	id := -1
	err := db.QueryRow(checkThreadExistAndGetIDBySlug, slug).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return id, false, nil
		}
		return id, false, err
	}
	return id, true, nil
}

func isThreadExist(db *pgx.ConnPool, id int) (bool, error) {
	err := db.QueryRow(checkThreadExistByID, id).Scan()
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ifThreadExistAndGetFodumSlugByID(db *pgx.ConnPool, id int) (string, bool, error) {
	forum := ""
	err := db.QueryRow(checkThreadExistAndGetForumSlugByID, id).Scan(&forum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forum, false, nil
		}
		return forum, false, err
	}
	return forum, true, nil
}

func ifThreadExistAndGetIDForumSlugBySlug(db *pgx.ConnPool, slug string) (string, int, bool, error) {
	id := -1
	forum := ""
	err := db.QueryRow(checkThreadExistAndGetIDForumSlugBySlug, slug).Scan(&id, &forum)
	if err != nil {
		if err == pgx.ErrNoRows {
			return forum, id, false, nil
		}
		return forum, id, false, err
	}
	return forum, id, true, nil
}

func dateTimeToString(date *strfmt.DateTime) pgx.NullString {
	time := pgx.NullString{}
	if date != nil {
		time.String = date.String()
		time.Valid = true
	} else {
		time.String = ""
		time.Valid = false
	}

	return time
}
