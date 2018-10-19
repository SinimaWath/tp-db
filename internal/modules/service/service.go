package service

import (
	"database/sql"

	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (pg ForumPgsql) Clear(operations.ClearParams) middleware.Responder {
	clear(pg.db)
	return nil
}

func clear(db *sql.DB) error {
	query := `TRUNCATE ONLY "user", forum, thread`
	_, err := db.Exec(query)
	return err
}
