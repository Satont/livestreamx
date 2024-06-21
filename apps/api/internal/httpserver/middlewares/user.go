package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

const USER_KEY = "user"

func (c *Middlewares) AttachUserToContext(ctx *gin.Context) {
	userID, err := c.sessionStorage.GetUserID(ctx)
	if err != nil {
		ctx.Next()
		return
	}
	if userID == "" {
		ctx.Next()
		return
	}

	user, _ := c.userRepo.FindByID(ctx.Request.Context(), uuid.MustParse(userID))
	if user == nil {
		ctx.Next()
		return
	}

	ctx.Set("user", user)
	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), USER_KEY, user))

	ctx.Next()
}

func GetUserFromContext(ctx context.Context) *user.User {
	user, ok := ctx.Value(USER_KEY).(*user.User)
	if !ok {
		return nil
	}
	return user
}
