package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self ForumPgsql) ForumCreate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ForumCreate")
	forum := &models.Forum{}
	err := forum.UnmarshalJSON(ctx.PostBody())
	if err != nil {
		log.Println(err)
		return
	}
	err = database.CreateForum(self.db, forum)
	if err != nil {
		switch err {
		case database.ErrForumNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		case database.ErrForumConflict:
			err := database.SelectForum(self.db, forum)
			if err != nil {
				log.Println("[ERROR] ForumCreate: " + err.Error())
				resp(ctx, Error, fasthttp.StatusInternalServerError)
				return
			}
			resp(ctx, forum, fasthttp.StatusConflict)
			return
		default:
			log.Println("[ERROR] ForumCreate: " + err.Error())
			resp(ctx, Error, fasthttp.StatusInternalServerError)
			return
		}
	}

	resp(ctx, forum, fasthttp.StatusCreated)
}

func (self *ForumPgsql) ForumGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ForumGetOne")
	forum := &models.Forum{}
	forum.Slug = ctx.UserValue("slug").(string)
	err := database.SelectForum(self.db, forum)
	if err != nil {
		if err == database.ErrForumNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		log.Println("[ERROR] ForumGetOne: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, forum, fasthttp.StatusOK)
}

func (self *ForumPgsql) ForumGetThreads(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ForumGetThreads")
	threads := &models.Threads{}

	err := database.SelectAllThreadsByForum(self.db, ctx.UserValue("slug").(string),
		ctx.QueryArgs().GetUintOrZero("limit"),
		getBool("desc", ctx.QueryArgs()),
		string(ctx.QueryArgs().Peek("since")), threads)

	if err != nil {
		if err == database.ErrForumNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, *threads, fasthttp.StatusOK)
	return
}

func (self *ForumPgsql) ForumGetUsers(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ForumGetUsers")
	users := &models.Users{}
	err := database.SelectAllUsersByForum(self.db, ctx.UserValue("slug").(string),
		ctx.QueryArgs().GetUintOrZero("limit"),
		getBool("desc", ctx.QueryArgs()),
		string(ctx.QueryArgs().Peek("since")), users)

	if err != nil {
		if err == database.ErrForumNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, *users, fasthttp.StatusOK)
	return
}
