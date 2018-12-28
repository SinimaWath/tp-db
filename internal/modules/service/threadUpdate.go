package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

func formThreadUpdateQuery(message, title, slug string, isID bool) (string, []interface{}) {
	query := `UPDATE thread SET `
	updateField := []string{}
	placeholder := 1
	resultArgs := []interface{}{}
	if message != "" {
		updateField = append(updateField, fmt.Sprintf("message = $%v", placeholder))
		resultArgs = append(resultArgs, message)
		placeholder++
	}

	if title != "" {
		updateField = append(updateField, fmt.Sprintf("title = $%v", placeholder))
		resultArgs = append(resultArgs, title)
		placeholder++
	}

	query += strings.Join(updateField, ", ")
	if isID {
		query += fmt.Sprintf(" WHERE id = $%v", placeholder)
	} else {
		query += fmt.Sprintf(" WHERE slug = $%v", placeholder)
	}

	// fmt.Println(query)
	resultArgs = append(resultArgs, slug)
	return query, resultArgs
}

func formThreadUpdateQueryLast(message, title string) (string, []interface{}) {
	query := `UPDATE thread SET `
	updateField := []string{}
	placeholder := 1
	resultArgs := []interface{}{}
	if message != "" {
		updateField = append(updateField, fmt.Sprintf("message = $%v", placeholder))
		resultArgs = append(resultArgs, message)
		placeholder++
	}

	if title != "" {
		updateField = append(updateField, fmt.Sprintf("title = $%v", placeholder))
		resultArgs = append(resultArgs, title)
		placeholder++
	}

	query += strings.Join(updateField, ", ")
	query += (" WHERE id IN (SELECT max(id) FROM thread)")
	return query, resultArgs
}

func checkIsAllNullString(nullString string, str ...string) bool {
	for _, s := range str {
		if s != nullString {
			return false
		}
	}

	return true
}

func (pg ForumPgsql) ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder {
	log.Println("ThreadUpdate")
	var queryUpdate string
	var args []interface{}

	if checkIsAllNullString("", params.Thread.Message, params.Thread.Title) {
		return pg.ThreadGetOne(operations.ThreadGetOneParams{
			SlugOrID: params.SlugOrID,
		})
	}

	if _, err := strconv.Atoi(params.SlugOrID); err != nil {
		queryUpdate, args = formThreadUpdateQuery(params.Thread.Message, params.Thread.Title, params.SlugOrID, false)
	} else {
		queryUpdate, args = formThreadUpdateQuery(params.Thread.Message, params.Thread.Title, params.SlugOrID, true)
	}

	result, err := pg.db.Exec(queryUpdate, args...)

	if err, ok := err.(*pq.Error); ok && err != nil {
		log.Println(err)
		return nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil
	} else if rowsAffected == 0 {
		return operations.NewThreadUpdateNotFound().WithPayload(&models.Error{})
	}

	return pg.ThreadGetOne(operations.ThreadGetOneParams{
		SlugOrID: params.SlugOrID,
	})
}
