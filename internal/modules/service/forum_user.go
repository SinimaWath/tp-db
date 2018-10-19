package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

func (pg ForumPgsql) UserCreate(params operations.UserCreateParams) middleware.Responder {

	queryInsert := "INSERT INTO \"user\" (nickname, fullname, about, email) VALUES ($1, $2, $3, $4);"
	_, err := pg.db.Exec(queryInsert, params.Nickname, params.Profile.Fullname, params.Profile.About, params.Profile.Email)

	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			log.Println(err)

			querySelect := `SELECT about, email, fullname, nickname FROM "user" WHERE nickname = $1 OR email = $2`
			rows, err := pg.db.Query(querySelect, params.Nickname, params.Profile.Email)

			if err != nil {
				log.Println(err)
				return nil
			}
			defer rows.Close()

			users := []*models.User{}
			for rows.Next() {
				user := &models.User{}
				scanErr := rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
				if scanErr != nil {
					return nil
				}
				users = append(users, user)
			}

			if rowsErr := rows.Err(); rowsErr != nil {
				log.Println(rowsErr)
				return nil
			}

			return operations.NewUserCreateConflict().WithPayload(users)
		}
	}
	params.Profile.Nickname = params.Nickname
	return operations.NewUserCreateCreated().WithPayload(params.Profile)
}

func (pg ForumPgsql) UserGetOne(params operations.UserGetOneParams) middleware.Responder {

	user := &models.User{}

	err := selectUser(pg.db, user, params.Nickname)
	switch err {
	case errNotFound:
		responseError := models.Error{"Can't find user"}
		return operations.NewUserGetOneNotFound().WithPayload(&responseError)
	case nil:
		return operations.NewUserGetOneOK().WithPayload(user)
	default:
		log.Println(err)
		return nil
	}
}
