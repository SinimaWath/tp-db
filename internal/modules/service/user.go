package service

import (
	"log"

	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) UserCreate(params operations.UserCreateParams) middleware.Responder {
	log.Println("[INFO] UserCreate")
	params.Profile.Nickname = params.Nickname

	err := database.CreateUser(self.db, params.Profile)
	if err != nil {
		switch err {
		case database.ErrUserConflict:
			users, err := database.SelectUsersWithNickOrEmail(self.db, params.Profile.Nickname, params.Profile.Email)
			if err != nil {
				log.Println("[ERROR] UserCreate: " + err.Error())
				return nil
			}
			return operations.NewUserCreateConflict().WithPayload(users)
		default:
			log.Println("[ERROR] UserCreate: " + err.Error())
			return nil
		}
	}

	return operations.NewUserCreateCreated().WithPayload(params.Profile)
}

func (self *ForumPgsql) UserGetOne(params operations.UserGetOneParams) middleware.Responder {
	log.Println("[INFO] UserGetOne")
	user := &models.User{}
	user.Nickname = params.Nickname
	err := database.SelectUser(self.db, user)

	if err != nil {
		if err == database.ErrUserNotFound {
			return operations.NewUserGetOneNotFound().WithPayload(&models.Error{})
		}
		log.Println("[ERROR] UserGetOne: " + err.Error())
		return nil
	}

	return operations.NewUserGetOneOK().WithPayload(user)
}

func (self *ForumPgsql) UserUpdate(params operations.UserUpdateParams) middleware.Responder {
	log.Println("[INFO] UserUpdate")
	user := &models.User{}
	user.Nickname = params.Nickname
	err := database.UpdateUser(self.db, user, params.Profile)
	if err != nil {
		switch err {
		case database.ErrUserNotFound:
			return operations.NewUserUpdateNotFound().WithPayload(&models.Error{})
		case database.ErrUserConflict:
			return operations.NewUserUpdateConflict().WithPayload(&models.Error{})
		}

		log.Println("[ERROR] UserUpdate: " + err.Error())
		return nil
	}

	return operations.NewUserUpdateOK().WithPayload(user)
}
