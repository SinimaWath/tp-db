package service

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

func (pg ForumPgsql) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	queryInsert := `INSERT INTO forum (user_nick, slug, title) VALUES ($1, $2, $3)`

	_, err := pg.db.Exec(queryInsert, params.Forum.User, params.Forum.Slug, params.Forum.Title)
	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			forum := &models.Forum{}
			if err := selectForum(pg.db, params.Forum.Slug, forum); err != nil {
				log.Println(err)
				return nil
			}
			return operations.NewForumCreateConflict().WithPayload(forum)
		}
		if err.Code == pgErrForeignKeyViolation {
			responseError := models.Error{Message: fmt.Sprintf("Can't find user with nickname: %v", params.Forum.User)}
			return operations.NewForumCreateNotFound().WithPayload(&responseError)
		}
		log.Println(err.Code)
		return nil
	}

	forum := &models.Forum{}
	if err := selectForum(pg.db, params.Forum.Slug, forum); err != nil {
		log.Println(err)
		return nil
	}
	return operations.NewForumCreateCreated().WithPayload(forum)
}

func (pg ForumPgsql) ForumGetOne(params operations.ForumGetOneParams) middleware.Responder {
	forum := &models.Forum{}
	err := selectForum(pg.db, params.Slug, forum)
	switch err {
	case errNotFound:
		responseError := models.Error{"Can't find user"}
		return operations.NewForumGetOneNotFound().WithPayload(&responseError)
	case nil:
		return operations.NewForumGetOneOK().WithPayload(forum)
	default:
		log.Println(err)
		return nil
	}
}

func selectForum(db *sql.DB, slug string, forum *models.Forum) error {
	querySelect := `SELECT u.nickname, f.slug, f.title FROM forum f JOIN "user" u ON u.nickname = f.user_nick WHERE slug = $1`

	row := db.QueryRow(querySelect, slug)

	if err := row.Scan(&forum.User, &forum.Slug, &forum.Title); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}
