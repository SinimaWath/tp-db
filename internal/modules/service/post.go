package service

import (
	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) PostsCreate(paramas operations.PostsCreateParams) middleware.Responder {
	posts, err := database.PostsCreate(self.db, paramas.SlugOrID, paramas.Posts)
	if err != nil {
		switch err {
		case database.ErrThreadNotFound, database.ErrUserNotFound:
			return operations.NewPostsCreateNotFound().WithPayload(&models.Error{})
		case database.ErrPostConflict:
			return operations.NewPostsCreateConflict().WithPayload(&models.Error{})
		}

		return nil
	}
	return operations.NewPostsCreateCreated().WithPayload(posts)
}

func (self *ForumPgsql) PostUpdate(params operations.PostUpdateParams) middleware.Responder {
	post := &models.Post{}
	post.ID = params.ID
	err := database.UpdatePost(self.db, post, params.Post)
	if err != nil {
		if err == database.ErrPostNotFound {
			return operations.NewPostUpdateNotFound().WithPayload(&models.Error{})
		}

		return nil
	}
	return operations.NewPostUpdateOK().WithPayload(post)
}

func (self *ForumPgsql) PostGetOne(params operations.PostGetOneParams) middleware.Responder {
	postFull := &models.PostFull{}
	postFull.Post = &models.Post{}

	postFull.Post.ID = params.ID
	err := database.SelectPostFull(self.db, params.Related, postFull)
	if err != nil {
		if err == database.ErrPostNotFound {
			return operations.NewPostGetOneNotFound().WithPayload(&models.Error{})
		}
		return nil
	}
	return operations.NewPostGetOneOK().WithPayload(postFull)
}
