package service

import (
	"database/sql"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

const (
	queryInsertForum      = `INSERT INTO forum (user_nick, slug, title) VALUES ($1, $2, $3)`
	queryCheckForumExists = `SELECT FROM forum where slug = $1`
	querySelectForum      = `SELECT u.nickname, f.slug, f.title FROM forum f JOIN "user" u ON u.nickname = f.user_nick WHERE slug = $1`
)

func (pg ForumPgsql) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	err := insertForum(pg.db, params.Forum.User, params.Forum.Slug, params.Forum.Title)

	if err == errForeignKeyViolation {
		return operations.NewForumCreateNotFound().WithPayload(&models.Error{})
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
	log.Println("ForumGetOne")
	forum := &models.Forum{}
	err := selectForumWithThreadsAndPosts(pg.db, params.Slug, forum)
	switch err {
	case errNotFound:
		return operations.NewForumGetOneNotFound().WithPayload(&models.Error{})
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

func (pg ForumPgsql) checkForumExist(slug string) bool {
	err := pg.db.QueryRow(queryCheckForumExists, slug).Scan()
	if err != nil {
		log.Println(slug + " checkForumExist ERROR: " + err.Error())
		return false
	}
	return true
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

const queryForumWithThreadsAndPosts = `
SELECT u.nickname, f.slug, f.title, f.post_count, f.thread_count
FROM forum f JOIN "user" u ON u.nickname = f.user_nick 
WHERE f.slug = $1`

func selectForumWithThreadsAndPosts(db *sql.DB, slug string, forum *models.Forum) error {
	row := db.QueryRow(queryForumWithThreadsAndPosts, slug)

	if err := row.Scan(&forum.User, &forum.Slug, &forum.Title, &forum.Posts, &forum.Threads); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}
