package service

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

const selectPostsFlatLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY p.created, p.id
	LIMIT $2
`

const selectPostsFlatLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY p.created DESC, p.id DESC
	LIMIT $2
`

const selectPostsFlatLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.id > $2
	ORDER BY p.created, p.id
	LIMIT $3
`
const selectPostsFlatLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.id < $2
	ORDER BY p.created DESC, p.id DESC
	LIMIT $3
`

const selectPostsTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY p.path
	LIMIT $2
`

const selectPostsTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1
	ORDER BY path DESC
	LIMIT $2
`

const selectPostsTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and (p.path > (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`

const selectPostsTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and (p.path < (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path DESC
	LIMIT $3
`

const selectPostsParentTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path
		LIMIT $3
	)
	ORDER BY path
`

const selectPostsParentTreeLimitDescByID = `
SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
FROM post p
WHERE p.thread_id = $1 and p.path[1] IN (
    SELECT p2.path[1]
    FROM post p2
	WHERE p2.parent_id IS NULL and p2.thread_id = $2
	ORDER BY p2.path DESC
    LIMIT $3
)
ORDER BY p.path[1] DESC, p.path[2:]
`

const selectPostsParentTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL and p2.path[1] > (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path
		LIMIT $4
	)
	ORDER BY p.path
`

const selectPostsParentTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, p.forum_slug
	FROM post p
	WHERE p.thread_id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread_id = $2 AND p2.parent_id IS NULL and p2.path[1] < (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path DESC
		LIMIT $4
	)
	ORDER BY p.path[1] DESC, p.path[2:]
`

const queryCheckThreadExistID = `
	SELECT true FROM thread where id = $1
`

const queryCheckThreadExistSlug = `
	SELECT true, id FROM thread where slug = $1
`

func (f *ForumPgsql) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	log.Println("ThreadGetPosts")
	isID := false
	if _, ok := strconv.Atoi(params.SlugOrID); ok == nil {
		isID = true
	}
	exist, id := checkThreadExistAndGetID(f.db, params.SlugOrID, isID)
	if !exist {
		return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{})
	}

	var rows *sql.Rows
	var selectErr error
	selectedPosts := models.Posts{}

	switch *params.Sort {
	case "flat":
		if params.Since != nil {
			if params.Desc != nil && *params.Desc == true {
				rows, selectErr = f.db.Query(selectPostsFlatLimitSinceDescByID, id,
					params.Since, params.Limit)
			} else {
				rows, selectErr = f.db.Query(selectPostsFlatLimitSinceByID, id,
					params.Since, params.Limit)
			}
		} else {
			if params.Desc != nil && *params.Desc == true {
				rows, selectErr = f.db.Query(selectPostsFlatLimitDescByID, id, params.Limit)
			} else {
				rows, selectErr = f.db.Query(selectPostsFlatLimitByID, id, params.Limit)
			}
		}
	case "tree":
		if params.Since != nil {
			if params.Desc != nil && *params.Desc {
				rows, selectErr = f.db.Query(selectPostsTreeLimitSinceDescByID, id,
					params.Since, params.Limit)
			} else {
				rows, selectErr = f.db.Query(selectPostsTreeLimitSinceByID, id,
					params.Since, params.Limit)
			}
		} else {
			if params.Desc != nil && *params.Desc {
				rows, selectErr = f.db.Query(selectPostsTreeLimitDescByID, id, params.Limit)
			} else {
				rows, selectErr = f.db.Query(selectPostsTreeLimitByID, id, params.Limit)
			}
		}
	case "parent_tree":
		if params.Since != nil {
			if params.Desc != nil && *params.Desc {
				rows, selectErr = f.db.Query(selectPostsParentTreeLimitSinceDescByID, id, id,
					params.Since, params.Limit)
			} else {
				rows, selectErr = f.db.Query(selectPostsParentTreeLimitSinceByID, id, id,
					params.Since, params.Limit)
			}
		} else {
			if params.Desc != nil && *params.Desc {
				rows, selectErr = f.db.Query(selectPostsParentTreeLimitDescByID, id, id,
					params.Limit)
			} else {
				rows, selectErr = f.db.Query(selectPostsParentTreeLimitByID, id, id,
					params.Limit)
			}
		}
	}

	if selectErr != nil {
		log.Println(selectErr)
		return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{})
	}

	defer rows.Close()
	for rows.Next() {
		post := &models.Post{}
		parentID := sql.NullInt64{}
		err := rows.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited,
			&post.Message, &parentID, &post.Thread, &post.Forum)

		if err != nil {
			log.Println(err)
			return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{})
		}

		if parentID.Valid {
			post.Parent = parentID.Int64
		} else {
			post.Parent = 0
		}

		selectedPosts = append(selectedPosts, post)
	}
	return operations.NewThreadGetPostsOK().WithPayload(selectedPosts)
}

func printPosts(posts *models.Posts) {

}
