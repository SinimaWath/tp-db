package service

import (
	"github.com/SinimaWath/tp-db/internal/models"
	"github.com/SinimaWath/tp-db/internal/modules/database"
	"github.com/SinimaWath/tp-db/internal/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func (self *ForumPgsql) UserCreate(params operations.UserCreateParams) middleware.Responder {
	params.Profile.Nickname = params.Nickname

	err := database.CreateUser(self.db, params.Profile)
	if err != nil {
		switch err {
		case database.ErrUserConflict:
			users, err := database.SelectUsersWithNickOrEmail(self.db, params.Profile.Nickname, params.Profile.Email)
			if err != nil {
				return nil
			}
			return operations.NewUserCreateConflict().WithPayload(users)
		default:
			return nil
		}
	}

	return operations.NewUserCreateCreated().WithPayload(params.Profile)
}

func (self *ForumPgsql) UserGetOne(params operations.UserGetOneParams) middleware.Responder {
	user := &models.User{}
	user.Nickname = params.Nickname
	err := database.SelectUser(self.db, user)

	if err != nil {
		if err == database.ErrUserNotFound {
			return operations.NewUserGetOneNotFound().WithPayload(&models.Error{})
		}
		return nil
	}

	return operations.NewUserGetOneOK().WithPayload(user)
}

func (self *ForumPgsql) UserUpdate(params operations.UserUpdateParams) middleware.Responder {
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

		return nil
	}

	return operations.NewUserUpdateOK().WithPayload(user)
}
