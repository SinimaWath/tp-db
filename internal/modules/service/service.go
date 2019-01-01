package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) Clear(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] Clear")
	err := database.Clear(self.db)
	if err != nil {
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	return
}

func (self *ForumPgsql) Status(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] Status")
	status := &models.Status{}
	err := database.Status(self.db, status)
	if err != nil {
		log.Println("[ERROR] Status: " + err.Error())
		return
	}
	resp(ctx, status, fasthttp.StatusOK)
	return
}
