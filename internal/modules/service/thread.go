package service

import (
	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) ThreadCreate(params operations.ThreadCreateParams) middleware.Responder {
	params.Thread.Forum = params.Slug
	err := database.ThreadCreate(self.db, params.Thread)

	if err != nil {
		switch err {
		case database.ErrThreadNotFoundAuthorOrForum:
			return operations.NewThreadCreateNotFound().WithPayload(&models.Error{})

		case database.ErrThreadConflict:
			return operations.NewThreadCreateConflict().WithPayload(params.Thread)
		}
		return nil
	}

	return operations.NewThreadCreateCreated().WithPayload(params.Thread)
}

func (self *ForumPgsql) ThreadGetOne(params operations.ThreadGetOneParams) middleware.Responder {
	thread := &models.Thread{}
	err := database.SelectThreadBySlugOrID(self.db, params.SlugOrID, thread)
	if err != nil {
		if err == database.ErrThreadNotFound {
			return operations.NewThreadGetOneNotFound().WithPayload(&models.Error{})
		}

		return nil
	}

	return operations.NewThreadGetOneOK().WithPayload(thread)
}

func (self *ForumPgsql) ThreadUpdate(params operations.ThreadUpdateParams) middleware.Responder {
	thread := &models.Thread{}
	err := database.UpdateThread(self.db, params.Thread, params.SlugOrID, thread)
	if err != nil {
		if err == database.ErrThreadNotFound {
			return operations.NewThreadUpdateNotFound().WithPayload(&models.Error{})
		}

		return nil
	}
	return operations.NewThreadUpdateOK().WithPayload(thread)
}

func (self *ForumPgsql) ThreadGetPosts(params operations.ThreadGetPostsParams) middleware.Responder {
	posts := &models.Posts{}
	err := database.SelectAllPostsByThread(self.db, params.SlugOrID, params.Limit,
		params.Desc, params.Since, params.Sort, posts)

	if err != nil {
		if err == database.ErrThreadNotFound {
			return operations.NewThreadGetPostsNotFound().WithPayload(&models.Error{})
		}

		return nil
	}
	return operations.NewThreadGetPostsOK().WithPayload(*posts)
}
