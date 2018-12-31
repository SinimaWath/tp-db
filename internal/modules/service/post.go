package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) PostsCreate(paramas operations.PostsCreateParams) middleware.Responder {
	log.Println("[INFO] PostsCreate")
	posts, err := database.PostsCreate(self.db, paramas.SlugOrID, paramas.Posts)
	if err != nil {
		switch err {
		case database.ErrThreadNotFound, database.ErrUserNotFound:
			return operations.NewPostsCreateNotFound().WithPayload(&models.Error{})
		case database.ErrPostConflict:
			return operations.NewPostsCreateConflict().WithPayload(&models.Error{})
		}

		log.Println("[ERROR] PostsCreate: " + err.Error())
		return nil
	}
	return operations.NewPostsCreateCreated().WithPayload(posts)
}

func (self *ForumPgsql) PostUpdate(params operations.PostUpdateParams) middleware.Responder {
	log.Println("[INFO] PostUpdate")
	post := &models.Post{}
	post.ID = params.ID
	err := database.UpdatePost(self.db, post, params.Post)
	if err != nil {
		if err == database.ErrPostNotFound {
			return operations.NewPostUpdateNotFound().WithPayload(&models.Error{})
		}

		log.Println("[ERROR] PostUpdate: " + err.Error())
		return nil
	}
	return operations.NewPostUpdateOK().WithPayload(post)
}

func (self *ForumPgsql) PostGetOne(params operations.PostGetOneParams) middleware.Responder {
	log.Println("[INFO] PostGetOne")
	postFull := &models.PostFull{}
	postFull.Post = &models.Post{}

	postFull.Post.ID = params.ID
	err := database.SelectPostFull(self.db, params.Related, postFull)
	if err != nil {
		if err == database.ErrPostNotFound {
			return operations.NewPostGetOneNotFound().WithPayload(&models.Error{})
		}
		log.Println("[ERROR] PostGetOne: " + err.Error())
		return nil
	}
	return operations.NewPostGetOneOK().WithPayload(postFull)
}
