package service

import (
	"fmt"
	"strconv"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

const queryUpdateCountID = `UPDATE thread t SET votes = (
	SELECT SUM(case when v.voice = true then 1 else -1 end)
	FROM vote v 
	WHERE v.thread_id=$1) WHERE id=$2`
const queryUpdateCountSLUG = `UPDATE thread SET votes = (
	SELECT SUM(case when v.voice = true then 1 else -1 end)
	FROM vote v
	JOIN thread t ON v.thread_id=t.id
	WHERE t.slug=$1
	) WHERE slug=$2`

const queryInsertVoteID = `INSERT INTO vote (nickname, voice, thread_id) VALUES ($1, $2, $3)
	ON CONFLICT ON CONSTRAINT unique_vote DO UPDATE SET voice = EXCLUDED.voice;`

const queryInsertVoteSLUG = `INSERT INTO vote (nickname, voice, thread_id) VALUES ($1, $2, 
	(select t.id from thread t where slug = $3))
	ON CONFLICT ON CONSTRAINT unique_vote DO UPDATE SET voice = EXCLUDED.voice;`

func voiceToBool(voice int32) bool {
	if voice == -1 {
		return false
	}
	return true
}

func (pg *ForumPgsql) ThreadVote(params operations.ThreadVoteParams) middleware.Responder {
	tx, err := pg.db.Begin()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	isID := false
	var execErr error
	if id, err := strconv.Atoi(params.SlugOrID); err != nil {
		_, execErr = tx.Exec(queryInsertVoteSLUG, params.Vote.Nickname, voiceToBool(params.Vote.Voice), params.SlugOrID)
	} else {
		isID = true
		_, execErr = tx.Exec(queryInsertVoteID, params.Vote.Nickname, voiceToBool(params.Vote.Voice), id)
	}

	if err, ok := execErr.(*pq.Error); ok && execErr != nil {
		tx.Rollback()
		fmt.Println(err)
		if err.Code == pgErrForeignKeyViolation {
			return operations.NewThreadVoteNotFound().WithPayload(&models.Error{})
		}
		return nil
	} else if execErr != nil {
		fmt.Println(execErr)
		tx.Rollback()
		return nil
	}

	if isID {
		_, execErr = tx.Exec(queryUpdateCountID, params.SlugOrID, params.SlugOrID)
	} else {
		_, execErr = tx.Exec(queryUpdateCountSLUG, params.SlugOrID, params.SlugOrID)
	}

	if err, ok := execErr.(*pq.Error); ok && execErr != nil {
		tx.Rollback()
		fmt.Println(err)
		if err.Code == pgErrForeignKeyViolation {
			return operations.NewThreadVoteNotFound().WithPayload(&models.Error{})
		}
		return nil
	} else if execErr != nil {
		fmt.Println(execErr)
		tx.Rollback()
		return nil
	}

	tx.Commit()

	thread := &models.Thread{}
	err = selectThreadVotes(pg.db, params.SlugOrID, isID, thread)
	if err != nil {
		return operations.NewThreadVoteNotFound().WithPayload(&models.Error{})
	}
	fmt.Println(thread.Votes)

	return operations.NewThreadUpdateOK().WithPayload(thread)
}
