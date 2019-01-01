package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) ThreadCreate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ThreadCreate")

	t := &models.Thread{}
	t.UnmarshalJSON(ctx.PostBody())
	t.Forum = ctx.UserValue("slug").(string)
	err := database.ThreadCreate(self.db, t)

	if err != nil {
		switch err {
		case database.ErrThreadNotFoundAuthorOrForum:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		case database.ErrThreadConflict:
			resp(ctx, t, fasthttp.StatusConflict)
			return
		}
		log.Println("[ERROR] ThreadCreate: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, t, fasthttp.StatusCreated)
	return
}

func (self *ForumPgsql) ThreadGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ThreadGetOne")
	thread := &models.Thread{}
	err := database.SelectThreadBySlugOrID(self.db, ctx.UserValue("slug_or_id").(string),
		thread)
	if err != nil {
		if err == database.ErrThreadNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, thread, fasthttp.StatusOK)
	return
}

func (self *ForumPgsql) ThreadUpdate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ThreadUpdate")

	thread := &models.Thread{}
	tU := &models.ThreadUpdate{}
	tU.UnmarshalJSON(ctx.PostBody())

	err := database.UpdateThread(self.db, tU, ctx.UserValue("slug_or_id").(string), thread)
	if err != nil {
		if err == database.ErrThreadNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}

		log.Println("[ERROR] ThreadUpdate: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, thread, fasthttp.StatusOK)
	return
}

func (self *ForumPgsql) ThreadGetPosts(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ThreadGetPosts")
	posts := &models.Posts{}

	err := database.SelectAllPostsByThread(self.db, ctx.UserValue("slug_or_id").(string),
		ctx.QueryArgs().GetUintOrZero("limit"), getBool("desc", ctx.QueryArgs()),
		ctx.QueryArgs().GetUintOrZero("since"),
		string(ctx.QueryArgs().Peek("sort")), posts)

	if err != nil {
		if err == database.ErrThreadNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}

		log.Println("[ERROR] ThreadGetPosts: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, *posts, fasthttp.StatusOK)
	return
}
