package database

import (
	"database/sql"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
)

const (
	updatePostFull = `
	UPDATE post SET message = $1
	WHERE id = $2
	RETURNING id, author, created, edited, message, parent_id, thread_id, forum_slug`
)

func UpdatePost(db *sql.DB, post *models.Post, pu *models.PostUpdate) error {
	var err error
	if pu.Message == "" {
		err = selectPost(db, post)
	} else {
		err = scanPost(db.QueryRow(updatePostFull, pu.Message, post.ID), post)
	}

	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return ErrPostNotFound
		}
		return err
	}
	return nil
}
