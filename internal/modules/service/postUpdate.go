package service

import (
	"database/sql"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

const updatePostQueryStart = `UPDATE post SET message = $1 where id = $2
RETURNING id, author, created, edited, message, parent_id, thread_id,
(select t.forum_slug from thread t where t.id = thread_id)`

func (f ForumPgsql) PostUpdate(params operations.PostUpdateParams) middleware.Responder {
	log.Println("PostUpdate")
	insertedPost := &models.Post{}

	if params.Post.Message == "" {
		err := selectPost(f.db, params.ID, insertedPost)
		if err != nil {
			log.Println(err)
			return nil
		}

		return operations.NewPostUpdateOK().WithPayload(insertedPost)
	}

	row := f.db.QueryRow(updatePostQueryStart, params.Post.Message, params.ID)
	parentId := sql.NullInt64{}
	err := row.Scan(&insertedPost.ID, &insertedPost.Author, &insertedPost.Created,
		&insertedPost.IsEdited, &insertedPost.Message, &parentId, &insertedPost.Thread,
		&insertedPost.Forum)

	if parentId.Valid {
		insertedPost.Parent = parentId.Int64
	} else {
		insertedPost.Parent = 0
	}

	if err != nil {
		log.Println(err)
		return operations.NewPostUpdateNotFound().WithPayload(&models.Error{})
	}

	return operations.NewPostUpdateOK().WithPayload(insertedPost)
}
