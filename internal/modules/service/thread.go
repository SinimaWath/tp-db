package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	log.Println("[INFO] ThreadCreate")
	params.Thread.Forum = params.Slug
	err := database.ThreadCreate(self.db, params.Thread)

	if err != nil {
		switch err {
		case database.ErrThreadNotFoundAuthorOrForum:
			return operations.NewThreadCreateNotFound().WithPayload(&models.Error{})

		case database.ErrThreadConflict:
			return operations.NewThreadCreateConflict().WithPayload(params.Thread)
		}
		log.Println("[ERROR] ThreadCreate: " + err.Error())
		return nil
	}

	return operations.NewThreadCreateCreated().WithPayload(params.Thread)
}

func (self *ForumPgsql) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	log.Println("[INFO] ThreadGetOne")
	thread := &models.Thread{}
	err := database.SelectThreadBySlugOrID(self.db, params.SlugOrID, thread)
	if err != nil {
		if err == database.ErrThreadNotFound {
			return operations.NewThreadGetOneNotFound().WithPayload(&models.Error{})
		}

		log.Println("[ERROR] ThreadGetOne: " + err.Error())
		return nil
	}

	return operations.NewThreadGetOneOK().WithPayload(thread)
}

func (self *ForumPgsql) ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder {
	log.Println("[INFO] ThreadUpdate")
	thread := &models.Thread{}
	err := database.UpdateThread(self.db, params.Thread, params.SlugOrID, thread)
	if err != nil {
		if err == database.ErrThreadNotFound {
			return operations.NewThreadUpdateNotFound().WithPayload(&models.Error{})
		}

		log.Println("[ERROR] ThreadUpdate: " + err.Error())
		return nil
	}
	return operations.NewThreadUpdateOK().WithPayload(thread)
}

func (self *ForumPgsql) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	log.Println("[INFO] ThreadGetPosts")
	posts := &models.Posts{}
	err := database.SelectAllPostsByThread(self.db, params.SlugOrID, params.Limit,
		params.Desc, params.Since, params.Sort, posts)

	if err != nil {
		if err == database.ErrThreadNotFound {
			return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{})
		}

		log.Println("[ERROR] ThreadGetPosts: " + err.Error())
		return nil
	}
	return operations.NewThreadGetPostsOK().WithPayload(*posts)
}
