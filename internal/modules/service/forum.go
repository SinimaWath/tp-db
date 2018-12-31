package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self ForumPgsql) ForumCreate(params operations.ForumCreateParams) middleware.Responder {
	log.Println("[INFO] ForumCreate")
	err := database.CreateForum(self.db, params.Forum)
	if err != nil {
		switch err {
		case database.ErrForumNotFound:
			return operations.NewForumCreateNotFound().WithPayload(&models.Error{})
		case database.ErrForumConflict:
			err := database.SelectForum(self.db, params.Forum)
			if err != nil {
				log.Println("[ERROR] ForumCreate: " + err.Error())
				return nil
			}
			return operations.NewForumCreateConflict().WithPayload(params.Forum)
		default:
			log.Println("[ERROR] ForumCreate: " + err.Error())
			return nil
		}
	}

	return operations.NewForumCreateCreated().WithPayload(params.Forum)
}

func (self *ForumPgsql) ForumGetOne(params operations.ForumGetOneParams) middleware.Responder {
	log.Println("[INFO] ForumGetOne")
	forum := &models.Forum{}
	forum.Slug = params.Slug
	err := database.SelectForum(self.db, forum)
	if err != nil {
		if err == database.ErrForumNotFound {
			return operations.NewForumGetOneNotFound().WithPayload(&models.Error{})
		}
		log.Println("[ERROR] ForumGetOne: " + err.Error())
		return nil
	}

	return operations.NewForumGetOneOK().WithPayload(forum)
}

func (self *ForumPgsql) ForumGetThreads(params operations.ForumGetThreadsParams) middleware.Responder {
	log.Println("[INFO] ForumGetThreads")
	threads := &models.Threads{}
	err := database.SelectAllThreadsByForum(self.db, params.Slug, params.Limit,
		params.Desc, params.Since, threads)

	if err != nil {
		if err == database.ErrForumNotFound {
			return operations.NewForumGetThreadsNotFound().WithPayload(&models.Error{})
		}
		log.Println("[ERROR] ForumGetThreads: " + err.Error())
		return nil
	}
	return operations.NewForumGetThreadsOK().WithPayload(*threads)
}

func (self *ForumPgsql) ForumGetUsers(params operations.ForumGetUsersParams) middleware.Responder {
	log.Println("[INFO] ForumGetUsers")
	users := &models.Users{}
	err := database.SelectAllUsersByForum(self.db, params.Slug, params.Limit,
		params.Desc, params.Since, users)

	if err != nil {
		if err == database.ErrForumNotFound {
			return operations.NewForumGetUsersNotFound().WithPayload(&models.Error{})
		}
		log.Println("[ERROR] ForumGetUsers: " + err.Error())
		return nil
	}
	return operations.NewForumGetUsersOK().WithPayload(*users)
}
