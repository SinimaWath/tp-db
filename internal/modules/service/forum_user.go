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
	querySelect := `SELECT about, email, fullname, nickname FROM "user" WHERE nickname = $1`

	row := pg.db.QueryRow(querySelect, params.Nickname)
	user := &models.User{}

	if err := row.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname); err != nil {
		if err == sql.ErrNoRows {
			responseError := models.Error{"Can't find user"}
			return operations.NewUserGetOneNotFound().WithPayload(&responseError)
		}
		log.Println(err)
		return nil
	}

	return operations.NewUserGetOneOK().WithPayload(user)
}

func userUpdateCheckNull(userUpdate *models.UserUpdate) bool {
	return userUpdate.Email == "" && userUpdate.About == "" && userUpdate.Fullname == ""
}

func formUserUpdateQuery(email, fullname, about, nickname string) (string, []interface{}) {
	query := `UPDATE "user" SET `
	updateField := []string{}
	placeholder := 1
	resultArgs := []interface{}{}
	if email != "" {
		updateField = append(updateField, fmt.Sprintf("email = $%v", placeholder))
		resultArgs = append(resultArgs, email)
		placeholder++
	}

	if fullname != "" {
		updateField = append(updateField, fmt.Sprintf("fullname = $%v", placeholder))
		resultArgs = append(resultArgs, fullname)
		placeholder++
	}
	if about != "" {
		updateField = append(updateField, fmt.Sprintf("about = $%v", placeholder))
		resultArgs = append(resultArgs, about)
		placeholder++
	}

	query += strings.Join(updateField, ", ")
	query += fmt.Sprintf(" WHERE nickname = $%v", placeholder)
	resultArgs = append(resultArgs, nickname)
	return query, resultArgs
}

func (pg ForumPgsql) UserUpdate(params operations.UserUpdateParams) middleware.Responder {

	if userUpdateCheckNull(params.Profile) {
		querySelect := `SELECT * from "user" WHERE nickname = $1`
		row := pg.db.QueryRow(querySelect, params.Nickname)
		user := &models.User{}
		if scanErr := row.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email); scanErr == sql.ErrNoRows {
			log.Println(sql.ErrNoRows)
			responseError := models.Error{"Can't find user"}
			return operations.NewUserUpdateNotFound().WithPayload(&responseError)
		}
		return operations.NewUserUpdateOK().WithPayload(user)
	}

	queryUpdate, args := formUserUpdateQuery(params.Profile.Email, params.Profile.Fullname, params.Profile.About, params.Nickname)
	log.Println(queryUpdate)
	log.Println(args)
	result, err := pg.db.Query(queryUpdate, args...)

	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			log.Println(err)
			responseError := &models.Error{"With such profile user already exist"}
			return operations.NewUserUpdateConflict().WithPayload(responseError)
		}
		log.Println(err)
		return nil

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil
	} else if rowsAffected == 0 {
		responseError := &models.Error{"Can't find user"}
		return operations.NewUserUpdateNotFound().WithPayload(responseError)
	}

	user := models.User{
		Nickname: params.Nickname,
		About:    params.Profile.About,
		Email:    params.Profile.Email,
		Fullname: params.Profile.Fullname,
	}
	return operations.NewUserUpdateOK().WithPayload(&user)
}
