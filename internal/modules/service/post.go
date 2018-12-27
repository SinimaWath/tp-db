package service

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	pq "github.com/lib/pq"
)

func (pg ForumPgsql) PostsCreate(params operations.PostsCreateParams) middleware.Responder {
	var insertErr error
	var insertedPosts *models.Posts
	if _, err := strconv.Atoi(params.SlugOrID); err == nil {
		insertedPosts, insertErr = postsInsert(pg.db, params.Posts, params.SlugOrID, true)
	} else {
		insertedPosts, insertErr = postsInsert(pg.db, params.Posts, params.SlugOrID, false)
	}

	if insertErr != nil {
		log.Println(insertErr)
		if err, ok := insertErr.(*pq.Error); ok {
			log.Printf("pqError: %#v", err)
			if err.Code == pgErrForeignKeyViolation {
				if err.Constraint == "post_parent_id_fkey" {
					return operations.NewPostsCreateConflict().WithPayload(&models.Error{})
				} else {
					return operations.NewPostsCreateNotFound().WithPayload(&models.Error{})
				}
			} else {
				if err.Column == "thread_id" {
					return operations.NewPostsCreateNotFound().WithPayload(&models.Error{})
				}
			}
			return nil
		}
		return nil
	}
	return operations.NewPostsCreateCreated().WithPayload(*insertedPosts)
}

const postInsertQuery = `
	INSERT INTO post (author, created, message, edited, parent_id, thread_id)
	VALUES 
`

const insertWithCheckParentID = "(SELECT (case when EXISTS(SELECT 1 from post p1 join thread t on p1.thread_id = t.id where p1.id=%v and t.id=%v) THEN %v ELSE -1 END))"
const insertWithCheckParentIDBySLUG = "(SELECT (case when EXISTS(SELECT 1 from post p1 join thread t on p1.thread_id = t.id where p1.id=%v and t.slug='%v') THEN %v ELSE -1 END))"

const postInsertEmptyToFindForeignKeyViolationsID = `INSERT INTO post (thread_id) VALUES ($1)`
const postInsertEmptyToFindForeignKeyViolationsSlug = `INSERT INTO post (thread_id) VALUES ((select id from thread where slug = $1))`

func formInsertValues(author, created, message, slugOrID string, isID bool, isEdited bool, parent int64, thread int32, placeholderStart int, valuesArgs *[]interface{}) string {
	values := "("
	valuesArr := []string{}
	placeholder := placeholderStart

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if author == "" {
		*valuesArgs = append(*valuesArgs, "NULL")
	} else {
		*valuesArgs = append(*valuesArgs, author)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if created == "" {
		*valuesArgs = append(*valuesArgs, "now()")
	} else {
		*valuesArgs = append(*valuesArgs, created)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	if message == "" {
		*valuesArgs = append(*valuesArgs, "NULL")
	} else {
		*valuesArgs = append(*valuesArgs, message)
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	placeholder++

	*valuesArgs = append(*valuesArgs, isEdited)

	if parent == 0 {
		valuesArr = append(valuesArr, fmt.Sprint("(NULL)"))
	} else {
		if isID {
			valuesArr = append(valuesArr, fmt.Sprintf(insertWithCheckParentID, parent, slugOrID, parent))
		} else {
			valuesArr = append(valuesArr, fmt.Sprintf(insertWithCheckParentIDBySLUG, parent, slugOrID, parent))
		}
	}

	if isID {
		valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
		if thread == 0 {
			*valuesArgs = append(*valuesArgs, slugOrID)
		} else {
			*valuesArgs = append(*valuesArgs, thread)
		}
	} else {
		valuesArr = append(valuesArr, fmt.Sprintf("(SELECT id from thread where slug = '%v')", slugOrID))
	}

	values += strings.Join(valuesArr, ", ")
	values += ")"

	return values
}

func postsInsert(db *sql.DB, posts models.Posts, slugOrID string, isID bool) (*models.Posts, error) {
	valuesArgsById := make([]interface{}, 0, len(posts)*6)
	valuesArgsBySlug := make([]interface{}, 0, len(posts)*5)

	finalInsertValues := ""
	createdStr := ""
	insertedPosts := &models.Posts{}

	if len(posts) == 0 {
		var err error
		if isID {
			_, err = db.Exec(postInsertEmptyToFindForeignKeyViolationsID, slugOrID)
		} else {
			_, err = db.Exec(postInsertEmptyToFindForeignKeyViolationsSlug, slugOrID)
		}
		return insertedPosts, err
	}

	for idx, post := range posts {

		if post.Created != nil {
			createdStr = post.Created.String()
		} else {
			createdStr = ""
		}
		insertValues := ""
		if isID {
			insertValues = formInsertValues(post.Author, createdStr, post.Message, slugOrID, isID,
				post.IsEdited, post.Parent, post.Thread, idx*5+1, &valuesArgsById)
		} else {
			insertValues = formInsertValues(post.Author, createdStr, post.Message, slugOrID, isID,
				post.IsEdited, post.Parent, post.Thread, idx*4+1, &valuesArgsBySlug)
		}
		if idx != 0 {
			finalInsertValues += ",\n"
		}
		finalInsertValues += insertValues
	}

	finalQuery := postInsertQuery + finalInsertValues
	finalQuery += " RETURNING id, author, created, edited, message, parent_id, thread_id, "
	finalQuery += "(select t.forum_slug from thread t where t.id = thread_id)"
	//fmt.Println(finalQuery)
	var rows *sql.Rows
	var queryError error
	if isID {
		//fmt.Printf("Values by id: %#v\n", valuesArgsById)
		rows, queryError = db.Query(finalQuery, valuesArgsById...)
	} else {
		//fmt.Printf("Values by slug: %#v\n", valuesArgsBySlug)
		rows, queryError = db.Query(finalQuery, valuesArgsBySlug...)
	}

	if queryError != nil {
		return nil, queryError
	}

	defer rows.Close()

	for rows.Next() {
		post := &models.Post{}
		parentId := sql.NullInt64{}
		err := rows.Scan(&post.ID, &post.Author, &post.Created,
			&post.IsEdited, &post.Message, &parentId, &post.Thread, &post.Forum)
		if err != nil {
			log.Println(err)
			continue
		}

		if parentId.Valid {
			post.Parent = parentId.Int64
		} else {
			post.Parent = 0
		}

		*insertedPosts = append(*insertedPosts, post)

	}

	//log.Printf("Inserted posts: %#v\n", insertedPosts)
	return insertedPosts, nil
}

const querySelectPostByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent_id, p.thread_id, t.forum_slug 
	from post p
	join thread t on p.thread_id = t.id
	where p.id = $1
`

func selectPost(db *sql.DB, id int64, post *models.Post) error {
	row := db.QueryRow(querySelectPostByID, id)
	parent_id := sql.NullInt64{}
	err := row.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited,
		&post.Message, &parent_id, &post.Thread, &post.Forum)
	if parent_id.Valid {
		post.Parent = parent_id.Int64
	} else {
		post.Parent = 0
	}
	return err
}

func (f ForumPgsql) PostGetOne(params operations.PostGetOneParams) middleware.Responder {
	post := &models.PostFull{
		Post: &models.Post{},
	}
	err := selectPost(f.db, params.ID, post.Post)
	log.Printf("Select post: %v\n", *post.Post)
	if err != nil {
		log.Println(err)
		return operations.NewPostGetOneNotFound().WithPayload(&models.Error{})
	}

	for _, table := range params.Related {
		switch table {
		case "user":
			post.Author = &models.User{}
			err = selectUser(f.db, post.Author, post.Post.Author)
			if err != nil {
				log.Println(err)
				return operations.NewPostGetOneNotFound().WithPayload(&models.Error{})
			}
		case "forum":
			post.Forum = &models.Forum{}
			err = selectForumWithThreadsAndPosts(f.db, post.Post.Forum, post.Forum)
			if err != nil {
				log.Println(err)
				return operations.NewPostGetOneNotFound().WithPayload(&models.Error{})
			}
		case "thread":
			post.Thread = &models.Thread{}
			err = selectThread(f.db, strconv.Itoa(int(post.Post.Thread)), true, post.Thread)
			if err != nil {
				log.Println(err)
				return operations.NewPostGetOneNotFound().WithPayload(&models.Error{})
			}
		}
	}
	return operations.NewPostGetOneOK().WithPayload(post)
}
