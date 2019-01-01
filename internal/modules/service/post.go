package service

import (
	"log"
	"strings"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) PostsCreate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] PostsCreate")

	p := models.Posts{}
	p.UnmarshalJSON(ctx.PostBody())

	posts, err := database.PostsCreate(self.db, ctx.UserValue("slug_or_id").(string), p)
	if err != nil {
		switch err {
		case database.ErrThreadNotFound, database.ErrUserNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		case database.ErrPostConflict:
			resp(ctx, Error, fasthttp.StatusConflict)
			return
		}

		log.Println("[ERROR] PostsCreate: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, posts, fasthttp.StatusCreated)
	return
}

func (self *ForumPgsql) PostUpdate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] PostUpdate")
	post := &models.Post{}
	post.ID = int32(postIDToInt(ctx))

	pU := &models.PostUpdate{}
	pU.UnmarshalJSON(ctx.PostBody())
	err := database.UpdatePost(self.db, post, pU)
	if err != nil {
		if err == database.ErrPostNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}

		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, post, fasthttp.StatusOK)
	return
}

func (self *ForumPgsql) PostGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] PostGetOne")
	postFull := &models.PostFull{}
	postFull.Post = &models.Post{}

	postFull.Post.ID = int32(postIDToInt(ctx))
	related := ctx.QueryArgs().Peek("related")
	err := database.SelectPostFull(self.db, strings.Split(string(related), ","), postFull)
	if err != nil {
		if err == database.ErrPostNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, postFull, fasthttp.StatusOK)
	return
}
