package database

import "database/sql"

const (
	updateForumThreadCount = `
	UPDATE forum f SET thread_count = thread_count + 1
	WHERE f.slug = $1`

	updateForumPostCountByThreadID = `
	UPDATE forum f SET post_count = post_count + $1
	FROM thread t
	WHERE t.forum_slug = f.slug AND t.id = $2`
)

func forumUpdateThreadCount(tx *sql.Tx, forumSlug string) error {
	_, err := tx.Exec(updateForumThreadCount, forumSlug)
	return err
}

func forumUpdatePostCountByThreadID(tx *sql.Tx, threadID int, postsCount int) error {
	_, err := tx.Exec(updateForumPostCountByThreadID, postsCount, threadID)
	return err
}
