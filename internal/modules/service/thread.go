package service

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	pq "github.com/lib/pq"
)

func (pg ForumPgsql) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	slug := params.Slug
	if params.Slug != params.Thread.Slug && params.Thread.Slug != "" {
		slug = params.Thread.Slug
	}

	err := insertThread(pg.db, slug, params.Thread.Author, params.Thread.Title, params.Thread.Message,
		params.Thread.Forum, params.Thread.Created)

	if err, ok := err.(*pq.Error); ok && err != nil {

		if err.Code == pgErrCodeUniqueViolation {
			thread := &models.Thread{}
			if err := selectThread(pg.db, slug, false, thread); err != nil {
				log.Println(err)
				return nil
			}
			// !!! для тестов
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
	if err := selectThread(pg.db, slug, false, thread); err != nil {
		log.Println(err)
		return nil
	}

	if params.Thread.Slug == "" {
		thread.Slug = ""
	}
	// !!! Для тестов
	return operations.NewThreadCreateCreated().WithPayload(thread)
}

func insertThread(db *sql.DB, slug, author, title, message, forum string, created *strfmt.DateTime) error {
	queryInsert := `INSERT INTO thread (slug, user_nick, created, forum_slug, title, message) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(queryInsert, slug, author, created,
		forum, title, message)

	if err != nil {
		return err
	}

	return nil
}

func selectThread(db *sql.DB, slugOrID string, isID bool, thread *models.Thread) error {
	querySelect := `SELECT f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id FROM thread t
	JOIN forum f ON f.slug = t.forum_slug `

	if isID {
		querySelect += `WHERE t.id = $1`
	} else {
		querySelect += `WHERE t.slug = $1`
	}

	row := db.QueryRow(querySelect, slugOrID)
	// thread.ID = 42
	if err := row.Scan(&thread.Forum, &thread.Author, &thread.Created, &thread.Slug,
		&thread.Title, &thread.Message, &thread.ID); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}

func selectThreadVotes(db *sql.DB, slugOrID string, isID bool, thread *models.Thread) error {
	querySelect := `SELECT f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id, t.votes FROM thread t
	JOIN forum f ON f.slug = t.forum_slug `

	if isID {
		querySelect += `WHERE t.id = $1`
	} else {
		querySelect += `WHERE t.slug = $1`
	}

	row := db.QueryRow(querySelect, slugOrID)
	// thread.ID = 42
	if err := row.Scan(&thread.Forum, &thread.Author, &thread.Created, &thread.Slug,
		&thread.Title, &thread.Message, &thread.ID, &thread.Votes); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}

func selectThreadID(db *sql.DB, id int, thread *models.Thread) error {
	querySelect := `SELECT t.id, f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id FROM thread t
	JOIN forum f ON f.slug = t.forum_slug
	WHERE t.id = $1`
	row := db.QueryRow(querySelect, id)
	// thread.ID = 42
	if err := row.Scan(&thread.ID, &thread.Forum, &thread.Author, &thread.Created, &thread.Slug,
		&thread.Title, &thread.Message, &thread.ID); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}

func selectThreads(db *sql.DB, slug, since string, limit int, desc bool, threads *models.Threads) error {
	query := `SELECT t.id, f.slug, u.nickname, t.created as created, t.slug, t.title, t.message FROM thread t
	JOIN forum f ON f.slug = t.forum_slug
	JOIN "user" u ON u.nickname = t.user_nick
	WHERE f.slug = $1`

	args := []interface{}{slug}
	placeholder := 2
	if since != "" {
		if desc {
			query += fmt.Sprintf(" AND t.created <= $%v", placeholder)
		} else {
			query += fmt.Sprintf(" AND t.created >= $%v", placeholder)
		}
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

	// log.Printf("%v, args: %#v\n", query, args)
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		thread := &models.Thread{}
		err := rows.Scan(&thread.ID, &thread.Forum, &thread.Author, &thread.Created, &thread.Slug, &thread.Title, &thread.Message)
		if err != nil {
			log.Println(err)
		}
		*threads = append(*threads, thread)
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

	if len(*threads) == 0 && !checkForumExist(pg.db, params.Slug) {
		responseError := models.Error{Message: fmt.Sprintf("Can't find forum by slug: %v", params.Slug)}
		return operations.NewForumGetThreadsNotFound().WithPayload(&responseError)
	} else if err != nil {
		log.Println(err)
		return nil
	}
	log.Printf("%#v\n", threads)
	return operations.NewForumGetThreadsOK().WithPayload(*threads)
}

func (pg ForumPgsql) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	thread := &models.Thread{}
	var selectErr error

	if _, err := strconv.Atoi(params.SlugOrID); err != nil {
		selectErr = selectThreadVotes(pg.db, params.SlugOrID, false, thread)
	} else {
		selectErr = selectThreadVotes(pg.db, params.SlugOrID, true, thread)
	}

	switch selectErr {
	case errNotFound:
		responseError := &models.Error{"Can't find thread"}
		return operations.NewThreadGetOneNotFound().WithPayload(responseError)
	case nil:
		return operations.NewThreadGetOneOK().WithPayload(thread)
	default:
		log.Println(selectErr)
		return nil
	}
}
