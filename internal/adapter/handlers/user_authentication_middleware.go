package handlers

import (
	"log/slog"
	"music-streaming/internal/core/ports"

	"github.com/gin-gonic/gin"
)

const RequestingUserKey = "requesting-user"

type UserManagementMiddleware struct {
	userAuthService ports.UserAuthenticationPort
	logger          *slog.Logger
}

func NewUserManagementMiddleware(userAuthServ ports.UserAuthenticationPort, logger *slog.Logger) *UserManagementMiddleware {
	return &UserManagementMiddleware{
		userAuthService: userAuthServ,
		logger:          logger,
	}
}

func (m *UserManagementMiddleware) WithAuth(c *gin.Context) {
	var (
		requiredParams  = c.MustGet(RequiredParameterKey).(requiredParams)
		qUser           = requiredParams.U
		qHashedPassword = requiredParams.T
		qSalt           = requiredParams.S
		ctx             = c.Request.Context()
	)

	m.logger.Info("Authentication middleware", slog.String("username", qUser))
	user, err := m.userAuthService.AuthenticateUser(ctx, qUser, qHashedPassword, qSalt)
	if err != nil {
		m.logger.Warn("Authentication failed", slog.String("username", qUser), slog.String("error", err.Error()))
		handleServiceError(c, err)
		return
	}

	m.logger.Info("Authentication successful", slog.String("username", qUser))
	c.Set(RequestingUserKey, user)
}
