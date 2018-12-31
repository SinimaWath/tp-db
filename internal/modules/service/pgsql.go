package service

import (
	"log"

	"gopkg.in/jackc/pgx.v2"
)

const postgres = "postgres"

const (
	pgErrCodeUniqueViolation = "23505"
	pgErrForeignKeyViolation = "23503"
)

type ForumPgsql struct {
	db *pgx.ConnPool
}

func NewForumPgsql(config *pgx.ConnConfig) *ForumPgsql {
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     *config,
		MaxConnections: 8,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	}
	log.Println(*config)
	p, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &ForumPgsql{
		db: p,
	}
}
