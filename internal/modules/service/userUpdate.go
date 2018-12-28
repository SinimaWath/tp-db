package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

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
	log.Println("UserUpdate")
	if userUpdateCheckNull(params.Profile) {
		return pg.UserGetOne(operations.UserGetOneParams{
			Nickname: params.Nickname,
		})
	}

	queryUpdate, args := formUserUpdateQuery(params.Profile.Email, params.Profile.Fullname, params.Profile.About, params.Nickname)
	result, err := pg.db.Exec(queryUpdate, args...)

	if err, ok := err.(*pq.Error); ok && err != nil {
		if err.Code == pgErrCodeUniqueViolation {
			log.Println(err)
			return operations.NewUserUpdateConflict().WithPayload(&models.Error{})
		}
		log.Println(err)
		return nil

	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil
	} else if rowsAffected == 0 {
		return operations.NewUserUpdateNotFound().WithPayload(&models.Error{})
	}

	return pg.UserGetOne(operations.UserGetOneParams{
		Nickname: params.Nickname,
	})
}
