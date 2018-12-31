package database

import (
	"database/sql"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/lib/pq"
)

const (
	insertThread = `
	INSERT INTO thread (slug, user_nick, created, forum_slug, title, message) 
	VALUES ($1,
	(SELECT u.nickname FROM "user" u WHERE u.nickname = $2),
	$3,
	(SELECT f.slug FROM forum f WHERE f.slug = $4)	
	,$5, $6)
	RETURNING id, slug, user_nick, created, forum_slug, title, message, votes
	`
)

func ThreadCreate(db *sql.DB, thread *models.Thread) error {

	tx, err := db.Begin()
	if err != nil {
		log.Println("[ERROR] ThreadCreate db.Begin(): " + err.Error())
		return err
	}

	err = scanThread(tx.QueryRow(insertThread, slugToNullable(thread.Slug), thread.Author, thread.Created, thread.Forum,
		thread.Title, thread.Message), thread)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] ThreadCreate tx.Rollback(): " + txErr.Error())
			return txErr
		}
		if err, ok := err.(*pq.Error); ok {
			switch err.Code {
			case pgErrCodeNotNullViolation, pgErrForeignKeyViolation:
				return ErrThreadNotFoundAuthorOrForum
			case pgErrCodeUniqueViolation:
				err := SelectThreadBySlug(db, thread)
				if err != nil {
					log.Println("[ERROR] ThreadCreate pgErrCodeUniqueViolation: " + err.Error())
					return err
				}

				return ErrThreadConflict
			}
		}
		return err
	}

	err = forumUpdateThreadCount(tx, thread.Forum)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] ThreadCreate tx.Rollback(): " + txErr.Error())
			return txErr
		}
		return err
	}

	err = createForumUserTx(tx, thread.Author, thread.Forum)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] ThreadCreate tx.Rollback(): " + txErr.Error())
			return txErr
		}
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Println("[ERROR] ThreadCreate tx.Commit(): " + commitErr.Error())
		return commitErr
	}

	return nil
}