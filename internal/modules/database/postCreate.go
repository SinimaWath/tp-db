package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/lib/pq"
)

func PostsCreate(db *sql.DB, slugOrIDThread string, posts models.Posts) (models.Posts, error) {
	isExist := false
	var err error
	threadID := 0
	if id, isID := isID(slugOrIDThread); isID {
		threadID = id
		isExist, err = isThreadExist(db, threadID)
	} else {
		threadID, isExist, err = ifThreadExistGetID(db, slugOrIDThread)
	}
	if !isExist {
		return nil, ErrThreadNotFound
	}

	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("[ERROR] PostsCreate db.Begin(): " + err.Error())
		return nil, err
	}

	resultPosts, err := insertPostsTx(tx, threadID, posts)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] PostsCreate tx.Rollback(): " + txErr.Error())
			return nil, txErr
		}

		if pqError, ok := err.(*pq.Error); ok && pqError != nil {
			switch pqError.Code {
			case pgErrForeignKeyViolation:
				if pqError.Constraint == "post_parent_id_fkey" {
					return nil, ErrPostConflict
				}
				if pqError.Constraint == "post_author_fkey" {
					return nil, ErrUserNotFound
				}
			}
		}
		return nil, err
	}

	err = forumUpdatePostCountByThreadID(tx, threadID, len(resultPosts))
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Println("[ERROR] PostsCreate tx.Rollback(): " + txErr.Error())
			return nil, txErr
		}
		return nil, err
	}
	if commitErr := tx.Commit(); commitErr != nil {
		log.Println("[ERROR] PostsCreate tx.Commit(): " + commitErr.Error())
		return nil, commitErr
	}

	return resultPosts, nil
}

const (
	insertPostsStart = `
	INSERT INTO post (forum_slug, author, created, message, edited, parent_id, thread_id)
	VALUES 
	`

	insertPostsEnd = `
	RETURNING id, author, created, edited, message, parent_id, thread_id, forum_slug
	`

	insertForumUsersStart = `
	INSERT INTO forum_user (nickname, forum_slug)
	VALUES
	`

	insertForumUsersEnd = `
	ON CONFLICT ON CONSTRAINT unique_forum_user DO NOTHING
	`
)

func insertPostsTx(tx *sql.Tx, threadID int, posts models.Posts) (models.Posts, error) {
	resultPosts := models.Posts{}
	if len(posts) == 0 {
		return resultPosts, nil
	}

	postsArgs := make([]interface{}, 0)
	forumUserArgs := make([]interface{}, 0)
	insertPostsQuery, insertForumUserQuery := formInsertQuery(threadID, posts, &postsArgs, &forumUserArgs)
	rows, queryError := tx.Query(*insertPostsQuery, postsArgs...)
	if queryError != nil {
		return nil, queryError
	}

	defer rows.Close()

	for rows.Next() {
		post := &models.Post{}
		err := scanPostRows(rows, post)
		if err != nil {
			return nil, err
		}

		resultPosts = append(resultPosts, post)
	}

	_, err := tx.Exec(*insertForumUserQuery, forumUserArgs...)

	return resultPosts, err
}

func formInsertQuery(id int, posts models.Posts,
	postsArgs *[]interface{}, forumUserArgs *[]interface{}) (*string, *string) {

	createdStr := ""
	insertValues := ""
	insertUserValues := ""
	finalInsertValues := strings.Builder{}
	finalUserInsertValues := strings.Builder{}

	for idx, post := range posts {
		if post.Created != nil {
			createdStr = post.Created.String()
		} else {
			createdStr = ""
		}

		insertValues = formInsertValuesID(post.Author, createdStr, post.Message, id,
			post.IsEdited, post.Parent, post.Thread, idx*5+1, postsArgs)

		insertUserValues = formInsertUserValues(post.Author, id, idx*2+1, forumUserArgs)

		if idx != 0 {
			finalInsertValues.WriteString(",")
			finalUserInsertValues.WriteString(",")
		}
		finalInsertValues.WriteString(insertValues)
		finalUserInsertValues.WriteString(insertUserValues)
	}

	insertPostsQuery := strings.Builder{}
	insertPostsQuery.WriteString(insertPostsStart)
	insertPostsQuery.WriteString(finalInsertValues.String())
	insertPostsQuery.WriteString(insertPostsEnd)

	insertUsersQuery := strings.Builder{}
	insertUsersQuery.WriteString(insertForumUsersStart)
	insertUsersQuery.WriteString(finalUserInsertValues.String())
	insertUsersQuery.WriteString(insertForumUsersEnd)

	resPosts, resUser := insertPostsQuery.String(), insertUsersQuery.String()
	return &resPosts, &resUser
}

func formInsertUserValues(author string, id, placeholder int, args *[]interface{}) string {
	values := fmt.Sprintf("($%v, (SELECT t.forum_slug FROM thread t WHERE t.id = $%v))", placeholder, placeholder+1)
	*args = append(*args, author)
	*args = append(*args, id)
	return values
}

const insertWithCheckParentID = `(
	SELECT (
		CASE WHEN 
		EXISTS(SELECT 1 from post p where p.id=%v and p.thread_id=%v)
		THEN %v ELSE -1 END)
	)`

func formInsertValuesID(author, created, message string, ID int, isEdited bool, parent int64, thread int32, placeholderStart int, valuesArgs *[]interface{}) string {
	values := "("
	valuesArr := []string{}
	placeholder := placeholderStart

	valuesArr = append(valuesArr, fmt.Sprintf(`(SELECT t.forum_slug from thread t where t.id = %v)`, ID))

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
		valuesArr = append(valuesArr, fmt.Sprintf(insertWithCheckParentID, parent, ID, parent))
	}

	valuesArr = append(valuesArr, fmt.Sprintf("$%v", placeholder))
	if thread == 0 {
		*valuesArgs = append(*valuesArgs, ID)
	} else {
		*valuesArgs = append(*valuesArgs, thread)
	}

	values += strings.Join(valuesArr, ", ")
	values += ")"

	return values
}
