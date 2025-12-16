package handlers

import (
	"music-streaming/internal/core/ports"

	"github.com/gin-gonic/gin"
)

const RequestingUserKey = "requesting-user"

type UserManagementMiddleware struct {
	userAuthService ports.UserAuthenticationPort
}

func NewUserManagementMiddleware(userAuthServ ports.UserAuthenticationPort) *UserManagementMiddleware {
	return &UserManagementMiddleware{
		userAuthService: userAuthServ,
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

	user, err := m.userAuthService.AuthenticateUser(ctx, qUser, qHashedPassword, qSalt)
	if err != nil {
		switch err.(type) {
		case *ports.FailedAuthenticationError:
			buildAndSendError(c, "40")
		}
		return
	}

	c.Set(RequestingUserKey, user)
}
