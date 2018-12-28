package service

import (
	"database/sql"
	"log"
)

const postgres = "postgres"

const (
	pgErrCodeUniqueViolation = "23505"
	pgErrForeignKeyViolation = "23503"
)

type ForumPgsql struct {
	db *sql.DB
}

func NewForumPgsql(dsn string) *ForumPgsql {
	db, err := sql.Open(postgres, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return &ForumPgsql{db: db}
}
