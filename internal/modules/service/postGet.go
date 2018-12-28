package service

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

const selectPostsFlatLimitBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1
	ORDER BY p.created, p.id
	LIMIT $2
`

const selectPostsFlatLimitDescBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1
	ORDER BY p.created DESC, p.id DESC
	LIMIT $2
`

const selectPostsFlatLimitSinceBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and p.id > $2
	ORDER BY p.created, p.id
	LIMIT $3
`
const selectPostsFlatLimitSinceDescBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and p.id < $2
	ORDER BY p.created DESC, p.id DESC
	LIMIT $3
`

const selectPostsFlatLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1
	ORDER BY p.created, p.id
	LIMIT $2
`

const selectPostsFlatLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1
	ORDER BY p.created DESC, p.id DESC
	LIMIT $2
`

const selectPostsFlatLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and p.id > $2
	ORDER BY p.created, p.id
	LIMIT $3
`
const selectPostsFlatLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and p.id < $2
	ORDER BY p.created DESC, p.id DESC
	LIMIT $3
`

const selectPostsTreeLimitBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1
	ORDER BY p.path
	LIMIT $2
`

const selectPostsTreeLimitDescBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1
	ORDER BY p.path DESC
	LIMIT $2
`

const selectPostsTreeLimitSinceBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and (p.path > (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`

const selectPostsTreeLimitSinceDescBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and (p.path < (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path DESC, p.id DESC
	LIMIT $3
`

const selectPostsTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1
	ORDER BY p.path
	LIMIT $2
`

const selectPostsTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1
	ORDER BY p.path DESC
	LIMIT $2
`

const selectPostsTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and (p.path > (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`

const selectPostsTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and (p.path < (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path DESC, p.id DESC
	LIMIT $3
`

const selectPostsParentTreeLimitBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.slug = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path[1]
		LIMIT $3
	)
	ORDER BY path
`

const selectPostsParentTreeLimitDescBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.slug = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path[1] DESC
		LIMIT $3 
	)
	ORDER BY p.path[1] DESC, p.path
`

const selectPostsParentTreeLimitSinceBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.slug = $2 AND p2.parent_id IS NULL and p2.path[1] > (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path[1]
		LIMIT $4
	)
	ORDER BY p.path
`

const selectPostsParentTreeLimitSinceDescBySlug = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.slug = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.slug = $2 AND p2.parent_id IS NULL and p2.path[1] < (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path[1] DESC
		LIMIT $4
	)
	ORDER BY p.path
`

const selectPostsParentTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.id = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path[1]
		LIMIT $3
	)
	ORDER BY path
`

const selectPostsParentTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.id = $2 AND p2.parent_id IS NULL
		ORDER BY p2.path[1] DESC
		LIMIT $3 
	)
	ORDER BY p.path[1] DESC, p.path
`

const selectPostsParentTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.id = $2 AND p2.parent_id IS NULL and p2.path[1] > (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path[1]
		LIMIT $4
	)
	ORDER BY p.path
`

const selectPostsParentTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug
	FROM post p
	JOIN thread t on t.id = p.thread_id
	WHERE t.id = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		JOIN thread t1 on p2.thread_id = t1.id
		WHERE t1.id = $2 AND p2.parent_id IS NULL and p2.path[1] < (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path[1] DESC
		LIMIT $4
	)
	ORDER BY p.path
`

const queryCheckThreadExistID = `
	SELECT true FROM thread where id = $1
`

const queryCheckThreadExistSlug = `
	SELECT true FROM thread where slug = $1
`

func checkThreadExist(db *sql.DB, slugOrId string, isID bool) bool {
	exist := sql.NullBool{}
	if isID {
		db.QueryRow(queryCheckThreadExistID, slugOrId).Scan(&exist)
	} else {
		db.QueryRow(queryCheckThreadExistSlug, slugOrId).Scan(&exist)
	}

	if exist.Valid {
		return true
	}

	return false
}

func (f *ForumPgsql) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	isID := false
	if _, ok := strconv.Atoi(params.SlugOrID); ok == nil {
		isID = true
	}

	if !checkThreadExist(f.db, params.SlugOrID, isID) {
		return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{})
	}

	var rows *sql.Rows
	var selectErr error
	selectedPosts := models.Posts{}
	if params.Desc != nil {
		log.Println("Desc: ", *params.Desc)
	}
	if params.Limit != nil {
		log.Println("Limit: ", *params.Limit)
	}
	if params.Since != nil {
		log.Println("Since: ", *params.Since)
	}
	if params.Sort != nil {
		log.Println("Sort: ", *params.Sort)
	}

	switch *params.Sort {
	case "flat":
		if isID {
			if params.Since != nil {
				if params.Desc != nil && *params.Desc == true {
					rows, selectErr = f.db.Query(selectPostsFlatLimitSinceDescByID, params.SlugOrID,
						params.Since, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsFlatLimitSinceByID, params.SlugOrID,
						params.Since, params.Limit)
				}
			} else {
				if params.Desc != nil && *params.Desc == true {
					rows, selectErr = f.db.Query(selectPostsFlatLimitDescByID, params.SlugOrID, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsFlatLimitByID, params.SlugOrID, params.Limit)
				}
			}
		} else {
			if params.Since != nil {
				if params.Desc != nil && *params.Desc == true {
					rows, selectErr = f.db.Query(selectPostsFlatLimitSinceDescBySlug, params.SlugOrID,
						params.Since, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsFlatLimitSinceBySlug, params.SlugOrID,
						params.Since, params.Limit)
				}
			} else {
				if params.Desc != nil && *params.Desc == true {
					rows, selectErr = f.db.Query(selectPostsFlatLimitDescBySlug, params.SlugOrID, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsFlatLimitBySlug, params.SlugOrID, params.Limit)
				}
			}
		}
	case "tree":
		if isID {
			if params.Since != nil {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsTreeLimitSinceDescByID, params.SlugOrID,
						params.Since, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsTreeLimitSinceByID, params.SlugOrID,
						params.Since, params.Limit)
				}
			} else {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsTreeLimitDescByID, params.SlugOrID, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsTreeLimitByID, params.SlugOrID, params.Limit)
				}
			}
		} else {
			if params.Since != nil {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsTreeLimitSinceDescBySlug, params.SlugOrID,
						params.Since, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsTreeLimitSinceBySlug, params.SlugOrID,
						params.Since, params.Limit)
				}
			} else {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsTreeLimitDescBySlug, params.SlugOrID, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsTreeLimitBySlug, params.SlugOrID, params.Limit)
				}
			}
		}
	case "parent_tree":
		if isID {
			if params.Since != nil {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitSinceDescByID, params.SlugOrID, params.SlugOrID,
						params.Since, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitSinceByID, params.SlugOrID, params.SlugOrID,
						params.Since, params.Limit)
				}
			} else {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitDescByID, params.SlugOrID, params.SlugOrID,
						params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitByID, params.SlugOrID, params.SlugOrID,
						params.Limit)
				}
			}
		} else {
			if params.Since != nil {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitSinceDescBySlug, params.SlugOrID, params.SlugOrID,
						params.Since, params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitSinceBySlug, params.SlugOrID, params.SlugOrID,
						params.Since, params.Limit)
				}
			} else {
				if params.Desc != nil && *params.Desc {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitDescBySlug, params.SlugOrID, params.SlugOrID,
						params.Limit)
				} else {
					rows, selectErr = f.db.Query(selectPostsParentTreeLimitBySlug, params.SlugOrID, params.SlugOrID,
						params.Limit)
				}
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
