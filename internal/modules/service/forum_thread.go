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

func formThreadCreateQuery(slug, author, title, message, created, forum string) (string, []interface{}) {
	query := `INSERT INTO thread (slug, user_nick, created, title, message`
	resultArgs := []interface{}{}
	resultArgs = append(resultArgs, []interface{}{slug, author, created, title, message}...)
	fPlaceHolder := ""
	if forum != "" {
		query += ", forum_slug"
		fPlaceHolder = ", $6"
		resultArgs = append(resultArgs, forum)
	}
	query += ") VALUES ($1, $2, $3, $4, $5"
	query += fPlaceHolder
	query += ")"

	return query, resultArgs
}

func (pg ForumPgsql) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	queryInsert := `INSERT INTO thread (slug, user_nick, created, forum_slug, title, message) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	slug := params.Slug
	if params.Slug != params.Thread.Slug && params.Thread.Slug != "" {
		slug = params.Thread.Slug
	}
	_, err := pg.db.Exec(queryInsert, slug, params.Thread.Author, params.Thread.Created,
		params.Thread.Forum, params.Thread.Title, params.Thread.Message)

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
			log.Println(err)
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

	if params.Thread.Slug == "" {
		thread.Slug = ""
	}
	return operations.NewThreadCreateCreated().WithPayload(thread)
}

func selectThread(db *sql.DB, slug string, thread *models.Thread) error {
	querySelect := `SELECT f.slug, t.user_nick, (t.created - interval '3 hour') AS created, t.slug, t.title, t.message FROM thread t
	JOIN forum f ON f.slug = t.forum_slug
	WHERE t.slug = $1`
	row := db.QueryRow(querySelect, slug)
	thread.ID = 42
	if err := row.Scan(&thread.Forum, &thread.Author, &thread.Created, &thread.Slug, &thread.Title, &thread.Message); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}

func selectThreads(db *sql.DB, slug, since string, limit int, desc bool, threads *models.Threads) error {
	query := `SELECT f.slug, u.nickname, (t.created - interval '3 hours') as created, t.slug, t.title, t.message FROM thread t
	JOIN forum f ON f.slug = t.forum_slug
	JOIN "user" u ON u.nickname = t.user_nick
	WHERE f.slug = $1`

	args := []interface{}{slug}
	placeholder := 2
	if since != "" {
		query += fmt.Sprintf(" AND t.created >= $%v", placeholder)
		placeholder++
		args = append(args, since)
	}

	query += " ORDER BY t.created"

	if desc != false {
		query += " DESC"
	}

	if limit != -1 {
		query += fmt.Sprintf(" LIMIT $%v", placeholder)
		args = append(args, limit)
	}

	log.Printf("%v, args: %#v\n", query, args)
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	len := 0
	defer rows.Close()
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(&thread.Forum, &thread.Author, &thread.Created, &thread.Slug, &thread.Title, &thread.Message)
		thread.ID = 42
		log.Printf("row: %#v\n", rows)
		if err != nil {
			log.Println(err)
		}
		*threads = append(*threads, thread)
		len++
	}

	log.Println(len)
	if len == 0 {
		return errNotFound
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (pg ForumPgsql) ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder {
	log.Println("ForumThread")
	threads := &models.Threads{}
	limit := -1
	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	since := ""
	if params.Since != nil {
		since = params.Since.String()
	}
	desc := false
	if params.Desc != nil {
		desc = *params.Desc
	}
	err := selectThreads(pg.db, params.Slug, since, limit, desc, threads)
	if err == errNotFound {
		responseError := models.Error{Message: fmt.Sprintf("Can't find forum by slug: %v", params.Slug)}
		return operations.NewForumGetThreadsNotFound().WithPayload(&responseError)
	} else if err != nil {
		log.Println(err)
		return nil
	}
	log.Printf("%#v\n", threads)
	return operations.NewForumGetThreadsOK().WithPayload(*threads)
}
