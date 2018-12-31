package database

import (
	"database/sql"
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/lib/pq"
)

const (
	insertVote = `
	INSERT INTO vote (nickname, voice, thread_id)
	VALUES ($1, $2, $3)
	ON CONFLICT ON CONSTRAINT unique_vote 
	DO UPDATE SET voice = EXCLUDED.voice;`
)

func VoteCreate(db *sql.DB, slugOrId string, t *models.Thread, v *models.Vote) error {

	if id, isID := isID(slugOrId); !isID {
		threadID, err := SelectThreadIDBySlug(db, slugOrId)
		if err != nil {
			return err
		}
		t.ID = int32(threadID)
	} else {
		t.ID = int32(id)
	}

	voteBool := voteIntToBool(v.Voice)
	tx, err := db.Begin()
	if err != nil {
		log.Println("[ERROR] VoteCreate db.Begin(): " + err.Error())
		return err
	}

	res, err := tx.Exec(insertVote, v.Nickname, voteBool, t.ID)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] VoteCreate tx.Rollback(): " + txErr.Error())
			return txErr
		}
		if pqError, ok := err.(*pq.Error); ok && pqError != nil {
			switch pqError.Code {
			case pgErrForeignKeyViolation:
				return ErrThreadNotFound
			}
		}
		return err
	}

	log.Println(res)

	err = threadUpdateVotesCountTx(tx, t)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] VoteCreate tx.Rollback(): " + txErr.Error())
			return txErr
		}
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Println("[ERROR] VoteCreate tx.Commit(): " + commitErr.Error())
		return commitErr
	}

	return nil
}
