package service

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

const (
	queryInsertUser           = "INSERT INTO \"user\" (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);"
	querySelectUser           = `SELECT about, email, fullname, nickname FROM "user" WHERE nickname = $1 OR email = $2`
	querySelectUserByNickname = `SELECT about, email, fullname, nickname FROM "user" WHERE nickname = $1`
)

func (pg ForumPgsql) UserCreate(params operations.UserCreateParams) middleware.Responder {

	err := insertUser(pg.db, params.Nickname, params.Profile.Fullname, params.Profile.About, params.Profile.Email)

	if err == errUniqueViolation {
		users := &models.Users{}
		selectErr := selectUsersByNicknameOrEmail(pg.db, users, params.Nickname, params.Profile.Email)
		if selectErr != nil {
			log.Println(selectErr)
			return nil
		}
		return operations.NewUserCreateConflict().WithPayload(*users)
	} else if err != nil {
		log.Println(err)
		return nil
	}

	params.Profile.Nickname = params.Nickname
	return operations.NewUserCreateCreated().WithPayload(params.Profile)
}

func (pg ForumPgsql) UserGetOne(params operations.UserGetOneParams) middleware.Responder {

	user := &models.User{}

	err := selectUser(pg.db, user, params.Nickname)
	switch err {
	case errNotFound:
		return operations.NewUserGetOneNotFound().WithPayload(&models.Error{})
	case nil:
		return operations.NewUserGetOneOK().WithPayload(user)
	default:
		log.Println(err)
		return nil
	}
}

func insertUser(db *sql.DB, nickname, fullname, about, email string) error {
	_, err := db.Exec(queryInsertUser, nickname, fullname, about, email)
	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			return errUniqueViolation
		}
		return err
	}
	return nil
}

func selectUsersByNicknameOrEmail(db *sql.DB, users *models.Users, nickname, email string) error {
	rows, err := db.Query(querySelectUser, nickname, email)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.User{}
		scanErr := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if scanErr != nil {
			return scanErr
		}
		*users = append(*users, user)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return rowsErr
	}
	return nil
}

func selectUser(db *sql.DB, user *models.User, nickname string) error {

	row := db.QueryRow(querySelectUserByNickname, nickname)

	if err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
		if err == sql.ErrNoRows {
			return errNotFound
		}
		return err
	}

	return nil
}

const querySelectStart = `
SELECT about, email, fullname, nickname from "user"
WHERE nickname IN (
    SELECT user_nick from thread where forum_slug = $1
    UNION
    SELECT p.author from post p
    JOIN thread t ON t.id = p.thread_id
    WHERE t.forum_slug = $1
)
`

func formForumGetUsersQuery(slug string, since *string, limit *int32, desc *bool) (string, []interface{}) {
	resultQuery := strings.Builder{}
	resultQuery.WriteString(querySelectStart)

	resultArgs := make([]interface{}, 0, 3)
	resultArgs = append(resultArgs, slug)

	placeholder := 2

	if since != nil {
		if desc != nil && *desc {
			resultQuery.WriteString(fmt.Sprintf(" AND nickname < $%v ", placeholder))
		} else {
			resultQuery.WriteString(fmt.Sprintf(" AND nickname > $%v ", placeholder))
		}
		resultArgs = append(resultArgs, *since)
		placeholder++
	}

	if desc != nil && *desc {
		resultQuery.WriteString("ORDER BY nickname DESC")
	} else {
		resultQuery.WriteString("ORDER BY nickname")
	}

	if limit != nil {
		resultQuery.WriteString(fmt.Sprintf("\n LIMIT $%v", placeholder))
		resultArgs = append(resultArgs, *limit)

	}

	return resultQuery.String(), resultArgs
}

func (pg ForumPgsql) ForumGetUsers(params operations.ForumGetUsersParams) middleware.Responder {
	if !checkForumExist(pg.db, params.Slug) {
		return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{})
	}

	query, args := formForumGetUsersQuery(params.Slug, params.Since, params.Limit, params.Desc)
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{})
	}
	defer rows.Close()

	users := models.Users{}

	for rows.Next() {
		user := &models.User{}
		scanErr := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if scanErr != nil {
			log.Println(scanErr)
			return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{})
		}
		users = append(users, user)
	}
	return operations.NewForumGetUsersOK().WithPayload(users)
}
