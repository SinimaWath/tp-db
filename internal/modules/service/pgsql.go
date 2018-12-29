package service

import (
	"database/sql"
	"log"
	"sync"
)

const postgres = "postgres"

const (
	pgErrCodeUniqueViolation = "23505"
	pgErrForeignKeyViolation = "23503"
)

type ForumPgsql struct {
	db *sql.DB

	sync.RWMutex
	forums map[string]struct{}
}

func NewForumPgsql(dsn string) *ForumPgsql {
	db, err := sql.Open(postgres, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return &ForumPgsql{
		db:     db,
		forums: make(map[string]struct{}, 20),
	}
}
