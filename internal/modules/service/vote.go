package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) ThreadVote(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] ThreadVote")
	thread := &models.Thread{}
	vote := &models.Vote{}

	vote.UnmarshalJSON(ctx.PostBody())
	err := database.VoteCreate(self.db, ctx.UserValue("slug_or_id").(string),
		thread, vote)

	if err != nil {
		switch err {
		case database.ErrThreadNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		log.Println("[ERROR] ThreadVote: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}
	resp(ctx, thread, fasthttp.StatusOK)
	return
}
