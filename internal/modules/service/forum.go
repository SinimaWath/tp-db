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

const (
	queryInsertForum      = `INSERT INTO forum (user_nick, slug, title) VALUES ($1, $2, $3)`
	queryCheckForumExists = `SELECT EXISTS(SELECT 1 FROM forum where slug = $1)`
	querySelectForum      = `SELECT u.nickname, f.slug, f.title FROM forum f JOIN "user" u ON u.nickname = f.user_nick WHERE slug = $1`
)

func (pg ForumPgsql) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	err := insertForum(pg.db, params.Forum.User, params.Forum.Slug, params.Forum.Title)

	if err == errForeignKeyViolation {
		responseError := models.Error{Message: fmt.Sprintf("Can't find user with nickname: %v", params.Forum.User)}
		return operations.NewForumCreateNotFound().WithPayload(&responseError)
	} else if err == errUniqueViolation {
		forum := &models.Forum{}
		if err := selectForum(pg.db, params.Forum.Slug, forum); err != nil {
			log.Println(err)
			return nil
		}
		return operations.NewForumCreateConflict().WithPayload(forum)
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
	err := selectForumWithThreadsAndPosts(pg.db, params.Slug, forum)
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

func insertForum(db *sql.DB, user, slug, title string) error {
	_, err := db.Exec(queryInsertForum, user, slug, title)
	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			return errUniqueViolation
		} else if err.Code == pgErrForeignKeyViolation {
			return errForeignKeyViolation
		}
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func checkForumExist(db *sql.DB, slug string) bool {
	row := db.QueryRow(queryCheckForumExists, slug)
	isExist := false
	row.Scan(&isExist)
	return isExist
}

func selectForum(db *sql.DB, slug string, forum *models.Forum) error {
	row := db.QueryRow(querySelectForum, slug)

	if err := row.Scan(&forum.User, &forum.Slug, &forum.Title); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}

const queryForumWithThreadsAndPosts = `SELECT u.nickname, f.slug, f.title,
 (select count(*) from thread t
 join post p on t.id = p.thread_id 
 where t.forum_slug = $1),
 (select count(*) from thread t 
 where t.forum_slug = $2)
FROM forum f JOIN "user" u ON u.nickname = f.user_nick WHERE slug = $3`

func selectForumWithThreadsAndPosts(db *sql.DB, slug string, forum *models.Forum) error {
	row := db.QueryRow(queryForumWithThreadsAndPosts, slug, slug, slug)

	if err := row.Scan(&forum.User, &forum.Slug, &forum.Title, &forum.Posts, &forum.Threads); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}
