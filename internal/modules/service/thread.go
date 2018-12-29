package service

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	pq "github.com/lib/pq"
)

const (
	queryInsertThread = `INSERT INTO thread (slug, user_nick, created, forum_slug, title, message) 
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	querySelectThreadByID = `SELECT f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id, t.votes FROM thread t
JOIN forum f ON f.slug = t.forum_slug WHERE t.id = $1`
	querySelectThreadBySlug = `SELECT f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id, t.votes FROM thread t
JOIN forum f ON f.slug = t.forum_slug WHERE t.slug = $1`

	querySelectThreadWithVotesByID = `SELECT f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id, t.votes FROM thread t
	JOIN forum f ON f.slug = t.forum_slug WHERE t.id = $1`
	querySelectThreadWithVotesBySlug = `SELECT f.slug, t.user_nick, t.created  AS created, t.slug, t.title, t.message, t.id, t.votes FROM thread t
	JOIN forum f ON f.slug = t.forum_slug WHERE t.slug = $1`

	querySelectThreadIDBySlug = `SELECT t.id from thread t where t.slug = $1`

	querySelectThreads = `SELECT t.id, f.slug, u.nickname, t.created as created, t.slug, t.title, t.message, t.votes FROM thread t
	JOIN forum f ON f.slug = t.forum_slug
	JOIN "user" u ON u.nickname = t.user_nick
	WHERE f.slug = $1`

	queryUpdateForumThreadCount = `
	UPDATE forum f SET thread_count = thread_count + 1
	WHERE f.slug = $1
	`
)

func (pg ForumPgsql) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	tx, err := pg.db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	id, err := insertThread(tx, params.Thread.Slug, params.Thread.Author, params.Thread.Title, params.Thread.Message,
		params.Thread.Forum, params.Thread.Created)

	if err, ok := err.(*pq.Error); ok && err != nil {
		tx.Rollback()
		if err.Code == pgErrCodeUniqueViolation {
			thread := &models.Thread{}
			if err := selectThread(pg.db, params.Thread.Slug, false, thread); err != nil {
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

	_, err = tx.Exec(queryUpdateForumThreadCount, params.Thread.Forum)

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil
	}
	if tx.Commit() != nil {
		log.Println(err)
		return nil
	}

	thread := &models.Thread{}
	if err := selectThread(pg.db, strconv.Itoa(id), true, thread); err != nil {
		log.Println(err)
		return nil
	}
	// !!! Для тестов
	return operations.NewThreadCreateCreated().WithPayload(thread)
}

func insertThread(db *sql.Tx, slug, author, title, message, forum string, created *strfmt.DateTime) (int, error) {

	nullableSlug := sql.NullString{}

	if slug == "" {
		nullableSlug.Valid = false
	} else {
		nullableSlug.String = slug
		nullableSlug.Valid = true
	}

	row := db.QueryRow(queryInsertThread, nullableSlug, author, created,
		forum, title, message)

	var id int

	err := row.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func selectThread(db *sql.DB, slugOrID string, isID bool, thread *models.Thread) error {
	var row *sql.Row
	if isID {
		row = db.QueryRow(querySelectThreadByID, slugOrID)
	} else {
		row = db.QueryRow(querySelectThreadBySlug, slugOrID)
	}

	nullableSlug := sql.NullString{}

	if err := row.Scan(&thread.Forum, &thread.Author, &thread.Created, &nullableSlug,
		&thread.Title, &thread.Message, &thread.ID, &thread.Votes); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	if nullableSlug.Valid {
		thread.Slug = nullableSlug.String
	} else {
		thread.Slug = ""
	}

	return nil
}

func selectThreadVotes(db *sql.DB, slugOrID string, isID bool, thread *models.Thread) error {
	var row *sql.Row
	if isID {
		row = db.QueryRow(querySelectThreadWithVotesByID, slugOrID)
	} else {
		row = db.QueryRow(querySelectThreadWithVotesBySlug, slugOrID)
	}

	nullableSlug := sql.NullString{}

	if err := row.Scan(&thread.Forum, &thread.Author, &thread.Created, &nullableSlug,
		&thread.Title, &thread.Message, &thread.ID, &thread.Votes); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	if nullableSlug.Valid {
		thread.Slug = nullableSlug.String
	} else {
		thread.Slug = ""
	}

	return nil
}

func selectThreadIDBySlug(db *sql.Tx, slug string) (int, error) {
	id := -1
	row := db.QueryRow(querySelectThreadIDBySlug, slug)
	err := row.Scan(&id)
	return id, err
}
func selectThreads(db *sql.DB, slug, since string, limit int, desc bool, threads *models.Threads) error {

	query := &strings.Builder{}
	query.WriteString(querySelectThreads)

	args := []interface{}{slug}
	placeholder := 2
	if since != "" {
		if desc {
			query.WriteString(fmt.Sprintf(" AND t.created <= $%v", placeholder))
		} else {
			query.WriteString(fmt.Sprintf(" AND t.created >= $%v", placeholder))
		}
		placeholder++
		args = append(args, since)
	}

	query.WriteString(" ORDER BY t.created")

	if desc != false {
		query.WriteString(" DESC")
	}

	if limit != -1 {
		query.WriteString(fmt.Sprintf(" LIMIT $%v", placeholder))
		args = append(args, limit)
	}

	rows, err := db.Query(query.String(), args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		thread := &models.Thread{}
		nullableSlug := sql.NullString{}
		err := rows.Scan(&thread.ID, &thread.Forum, &thread.Author, &thread.Created, &nullableSlug, &thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			log.Println(err)
		}
		if nullableSlug.Valid {
			thread.Slug = nullableSlug.String
		} else {
			thread.Slug = ""
		}

		*threads = append(*threads, thread)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func (pg ForumPgsql) ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder {
	log.Println("ForumGetThreads")
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

	if len(*threads) == 0 && !pg.checkForumExist(params.Slug) {
		return operations.NewForumGetThreadsNotFound().WithPayload(&models.Error{})
	} else if err != nil {
		log.Println("ForumGetThreads ERROR: " + err.Error())
		return nil
	}
	return operations.NewForumGetThreadsOK().WithPayload(*threads)
}

func (pg ForumPgsql) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	log.Println("ThreadGetOne")
	thread := &models.Thread{}
	var selectErr error

	if _, err := strconv.Atoi(params.SlugOrID); err != nil {
		selectErr = selectThreadVotes(pg.db, params.SlugOrID, false, thread)
	} else {
		selectErr = selectThreadVotes(pg.db, params.SlugOrID, true, thread)
	}

	switch selectErr {
	case errNotFound:
		return operations.NewThreadGetOneNotFound().WithPayload(&models.Error{})
	case nil:
		return operations.NewThreadGetOneOK().WithPayload(thread)
	default:
		log.Println(selectErr)
		return nil
	}
}

func checkThreadExistAndGetID(db *sql.DB, slugOrId string, isID bool) (bool, string) {
	exist := sql.NullBool{}
	id := ""
	if isID {
		db.QueryRow(queryCheckThreadExistID, slugOrId).Scan(&exist)
		id = slugOrId
	} else {
		db.QueryRow(queryCheckThreadExistSlug, slugOrId).Scan(&exist, &id)
	}

	if exist.Valid {
		return true, id
	}

	return false, ""
}
