package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

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
	log.Println("ForumGetUsers")
	if !pg.checkForumExist(params.Slug) {
		return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{})
	}

	query, args := formForumGetUsersQuery(params.Slug, params.Since, params.Limit, params.Desc)
	rows, err := pg.db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{})
	}

	users := models.Users{}

	defer rows.Close()
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
