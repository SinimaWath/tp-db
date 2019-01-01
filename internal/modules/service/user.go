package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/valyala/fasthttp"
)

func (self *ForumPgsql) UserCreate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] UserCreate")

	u := &models.User{}
	u.Nickname = ctx.UserValue("nickname").(string)
	u.UnmarshalJSON(ctx.PostBody())
	err := database.CreateUser(self.db, u)

	if err != nil {
		switch err {
		case database.ErrUserConflict:
			users, err := database.SelectUsersWithNickOrEmail(self.db, u.Nickname, u.Email)
			if err != nil {
				log.Println("[ERROR] UserCreate: " + err.Error())
				resp(ctx, Error, fasthttp.StatusInternalServerError)
				return
			}
			resp(ctx, users, fasthttp.StatusConflict)
			return
		default:
			log.Println("[ERROR] UserCreate: " + err.Error())
			resp(ctx, Error, fasthttp.StatusInternalServerError)
			return
		}
	}

	resp(ctx, u, fasthttp.StatusCreated)
	return
}

func (self *ForumPgsql) UserGetOne(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] UserGetOne")
	user := &models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	err := database.SelectUser(self.db, user)

	if err != nil {
		if err == database.ErrUserNotFound {
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		}
		log.Println("[ERROR] UserGetOne: " + err.Error())
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, user, fasthttp.StatusOK)
}

func (self *ForumPgsql) UserUpdate(ctx *fasthttp.RequestCtx) {
	log.Println("[INFO] UserUpdate")

	user := &models.User{}
	user.Nickname = ctx.UserValue("nickname").(string)
	userUpdate := &models.UserUpdate{}
	userUpdate.UnmarshalJSON(ctx.PostBody())

	err := database.UpdateUser(self.db, user, userUpdate)
	if err != nil {
		switch err {
		case database.ErrUserNotFound:
			resp(ctx, Error, fasthttp.StatusNotFound)
			return
		case database.ErrUserConflict:
			resp(ctx, Error, fasthttp.StatusConflict)
			return
		}
		resp(ctx, Error, fasthttp.StatusInternalServerError)
		return
	}

	resp(ctx, user, fasthttp.StatusOK)
	return
}
