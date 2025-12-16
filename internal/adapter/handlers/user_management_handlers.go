package handlers

import (
	"context"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserManagementHandler struct {
	userServ ports.UserManagementPort
}

func NewUserManagementHandler(userServ ports.UserManagementPort) *UserManagementHandler {
	return &UserManagementHandler{
		userServ: userServ,
	}
}

func (h *UserManagementHandler) hangleGetUser(c *gin.Context) {
	var (
		rUser    = c.MustGet(RequestingUserKey).(*domain.User)
		username = c.Query("username")
		ctx      = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	user, err := h.userServ.GetUser(ctx, username)
	if err != nil {
		switch err.(type) {
		case *ports.UserNotFoundError:
			buildAndSendError(c, "0")
		case *ports.UserNotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.MissingOrInvalidParameterError:
			buildAndSendError(c, "10")
		}
		return
	}

	user.Password = "" //do not send user passwords

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		User:    &user,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) hangleGetUsers(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	users, err := h.userServ.GetUsers(ctx)
	if err != nil {
		switch err.(type) {
		case *ports.UserNotAuthorizedError:
			buildAndSendError(c, "50")
		}
		return
	}

	for i := range users {
		users[i].Password = "" //do not send user passwords
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		Users:   &users,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) handleCreateUser(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	user := &domain.User{}
	if err := c.Bind(user); err != nil {
		buildAndSendError(c, "0")
		return
	}

	if err := h.userServ.CreateUser(ctx, *user); err != nil {
		switch err.(type) {
		case *ports.UserNotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.MissingOrInvalidParameterError:
			buildAndSendError(c, "10")
		default:
			buildAndSendError(c, "0")
		}
		return
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) handleUpdateUser(c *gin.Context) {
	var (
		rUser    = c.MustGet(RequestingUserKey).(*domain.User)
		username = c.PostForm("username")
		ctx      = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	user := &domain.User{}
	if err := c.Bind(user); err != nil {
		buildAndSendError(c, "0")
		return
	}

	//if no valid updates abort
	if err := h.userServ.UpdateUser(ctx, username, *user); err != nil {
		switch err.(type) {
		case *ports.UserNotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.MissingOrInvalidParameterError:
			buildAndSendError(c, "10")
		default:
			buildAndSendError(c, "0")
		}
		return
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) handleDeleteUser(c *gin.Context) {
	var (
		rUser    = c.MustGet(RequestingUserKey).(*domain.User)
		username = c.PostForm("username")
		ctx      = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	if err := h.userServ.DeleteUser(ctx, username); err != nil {
		switch err.(type) {
		case *ports.UserNotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.MissingOrInvalidParameterError:
			buildAndSendError(c, "10")
		default:
			buildAndSendError(c, "0")
		}
		return
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) handleChangePassword(c *gin.Context) {
	var (
		rUser    = c.MustGet(RequestingUserKey).(*domain.User)
		username = c.PostForm("username")
		password = c.PostForm("password")
		ctx      = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	if err := h.userServ.ChangePassword(ctx, username, password); err != nil {
		switch err.(type) {
		case *ports.UserNotAuthorizedError:
			buildAndSendError(c, "50")
		case *ports.MissingOrInvalidParameterError:
			buildAndSendError(c, "10")
		default:
			buildAndSendError(c, "0")
		}
		return
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
