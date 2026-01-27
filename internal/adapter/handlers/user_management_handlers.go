package handlers

import (
	"context"
	"log/slog"
	"music-streaming/internal/core/domain"
	"music-streaming/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserManagementHandler struct {
	userServ ports.UserManagementPort
	logger   *slog.Logger
}

func NewUserManagementHandler(userServ ports.UserManagementPort, logger *slog.Logger) *UserManagementHandler {
	return &UserManagementHandler{
		userServ: userServ,
		logger:   logger,
	}
}

func (h *UserManagementHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/getUser", h.hangleGetUser)
	group.GET("/getUsers", h.hangleGetUsers)
	group.POST("/createUser", h.handleCreateUser)
	group.POST("/updateUser", h.handleUpdateUser)
	group.POST("/deleteUser", h.handleDeleteUser)
	group.POST("/changePassword", h.handleChangePassword)
}

func (h *UserManagementHandler) hangleGetUser(c *gin.Context) {
	var (
		rUser    = c.MustGet(RequestingUserKey).(*domain.User)
		username = c.Query("username")
		ctx      = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	h.logger.Info("Get user handler called", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
	user, err := h.userServ.GetUser(ctx, username)
	if err != nil {
		h.logger.Warn("Get user handler error", slog.String("requesting_user", rUser.Username), slog.String("target_username", username), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	h.logger.Info("Get user handler success", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))

	// Convert to DTO and clear password
	userDTO := UserToDTO(user)
	userDTO.Password = "" // do not send user passwords

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		User:    &userDTO,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) hangleGetUsers(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	h.logger.Info("Get users handler called", slog.String("requesting_user", rUser.Username))
	users, err := h.userServ.GetUsers(ctx)
	if err != nil {
		h.logger.Warn("Get users handler error", slog.String("requesting_user", rUser.Username), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	h.logger.Info("Get users handler success", slog.String("requesting_user", rUser.Username), slog.Int("count", len(users)))

	// Convert to DTOs and clear passwords
	userDTOs := UsersToDTO(users)
	for i := range userDTOs {
		userDTOs[i].Password = "" // do not send user passwords
	}

	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
		Users:   &userDTOs,
	}

	SerializeAndSendBody(c, subsonicRes)
}

func (h *UserManagementHandler) handleCreateUser(c *gin.Context) {
	var (
		rUser = c.MustGet(RequestingUserKey).(*domain.User)
		ctx   = context.WithValue(c.Request.Context(), ports.KeyRequestingUserID, rUser)
	)

	userDTO := &UserDTO{}
	if err := c.Bind(userDTO); err != nil {
		h.logger.Warn("Create user handler - bind error", slog.String("requesting_user", rUser.Username), slog.String("error", err.Error()))
		buildAndSendError(c, "0")
		return
	}

	// Convert DTO to domain entity
	user := DTOToUser(*userDTO)

	h.logger.Info("Create user handler called", slog.String("requesting_user", rUser.Username), slog.String("target_username", user.Username))
	if err := h.userServ.CreateUser(ctx, user); err != nil {
		h.logger.Warn("Create user handler error", slog.String("requesting_user", rUser.Username), slog.String("target_username", user.Username), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	h.logger.Info("Create user handler success", slog.String("requesting_user", rUser.Username), slog.String("target_username", user.Username))
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

	userDTO := &UserDTO{}
	if err := c.Bind(userDTO); err != nil {
		h.logger.Warn("Update user handler - bind error", slog.String("requesting_user", rUser.Username), slog.String("error", err.Error()))
		buildAndSendError(c, "0")
		return
	}

	// Convert DTO to domain entity
	user := DTOToUser(*userDTO)

	h.logger.Info("Update user handler called", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
	if err := h.userServ.UpdateUser(ctx, username, user); err != nil {
		h.logger.Warn("Update user handler error", slog.String("requesting_user", rUser.Username), slog.String("target_username", username), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	h.logger.Info("Update user handler success", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
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

	h.logger.Info("Delete user handler called", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
	if err := h.userServ.DeleteUser(ctx, username); err != nil {
		h.logger.Warn("Delete user handler error", slog.String("requesting_user", rUser.Username), slog.String("target_username", username), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	h.logger.Info("Delete user handler success", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
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

	h.logger.Info("Change password handler called", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
	if err := h.userServ.ChangePassword(ctx, username, password); err != nil {
		h.logger.Warn("Change password handler error", slog.String("requesting_user", rUser.Username), slog.String("target_username", username), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	h.logger.Info("Change password handler success", slog.String("requesting_user", rUser.Username), slog.String("target_username", username))
	subsonicRes := SubsonicResponse{
		Xmlns:   Xmlns,
		Status:  "ok",
		Version: SubsonicVersion,
	}

	SerializeAndSendBody(c, subsonicRes)
}
