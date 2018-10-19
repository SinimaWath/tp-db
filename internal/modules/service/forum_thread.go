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

func (pg ForumPgsql) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	queryInsert := `INSERT INTO thread (forum_slug, user_nick, created, slug, title, message) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	slug := params.Slug
	if params.Slug != params.Thread.Slug && params.Thread.Slug != "" {
		slug = params.Thread.Slug
	}
	_, err := pg.db.Exec(queryInsert, params.Thread.Forum, params.Thread.Author, params.Thread.Created,
		slug, params.Thread.Title, params.Thread.Message)

	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			thread := &models.Thread{}
			if err := selectThread(pg.db, slug, thread); err != nil {
				log.Println(err)
				return nil
			}
			return operations.NewThreadCreateConflict().WithPayload(thread)
		}
		if err.Code == pgErrForeignKeyViolation {
			responseError := models.Error{Message: fmt.Sprintf("Can't find user with nickname: %v")}
			return operations.NewThreadCreateNotFound().WithPayload(&responseError)
		}
		log.Println(err.Code)
		return nil
	}

	thread := &models.Thread{}
	if err := selectThread(pg.db, slug, thread); err != nil {
		log.Println(err)
		return nil
	}
	thread.ID = 42

	if params.Thread.Slug == "" {
		thread.Slug = ""
	}

	log.Printf("%#v", thread)
	return operations.NewThreadCreateCreated().WithPayload(thread)
}

func selectThread(db *sql.DB, slug string, thread *models.Thread) error {
	querySelect := `SELECT forum_slug, user_nick, (created - interval '3 hour') AS created, slug, title, message FROM thread WHERE slug = $1`
	row := db.QueryRow(querySelect, slug)

	if err := row.Scan(&thread.Forum, &thread.Author, &thread.Created, &thread.Slug, &thread.Title, &thread.Message); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}
