package service

import (
	"database/sql"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (pg ForumPgsql) Clear(operations.ClearParams) middleware.Responder {
	log.Println("Clear")
	err := clear(pg.db)
	if err != nil {
		log.Println(err)
	}
	return nil
}

const queryStatus = `
SELECT COUNT(*) FROM forum
UNION ALL
SELECT COUNT(*) FROM thread 
UNION ALL
SELECT COUNT(*) FROM post
UNION ALL
SELECT COUNT(*) FROM "user"
`

func (f *ForumPgsql) Status(params operations.StatusParams) middleware.Responder {
	log.Println("Status")
	rows, err := f.db.Query(queryStatus)
	if err != nil {
		log.Println(err)
		return nil
	}

	status := models.Status{}
	rows.Next()
	rows.Scan(&status.Forum)
	rows.Next()
	rows.Scan(&status.Thread)
	rows.Next()
	rows.Scan(&status.Post)
	rows.Next()
	rows.Scan(&status.User)

	return operations.NewStatusOK().WithPayload(&status)
}

func clear(db *sql.DB) error {
	query := `TRUNCATE ONLY post, vote, thread, forum, "user"`
	_, err := db.Exec(query)
	return err
}
