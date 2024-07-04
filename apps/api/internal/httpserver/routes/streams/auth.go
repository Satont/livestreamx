package streams

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authReq struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Action   string `json:"action"`
	Path     string `json:"path"`
	Protocol string `json:"protocol"`
	Id       string `json:"id"`
	Query    string `json:"query"`
}

func (c *Streams) authHandler(ctx *gin.Context) {
	body := authReq{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(200)
	return

	if body.Action == "read" {
		ctx.Status(200)
		return
	}

	if body.Action == "publish" {
		streamKey, err := uuid.Parse(body.Path)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		user, err := c.userRepo.FindByStreamKey(ctx, streamKey)
		if err != nil {
			ctx.JSON(404, gin.H{"error": err.Error()})
			return
		}

		if user.Banned {
			ctx.JSON(403, gin.H{"error": "User is banned"})
			return
		}

		ctx.Status(200)
		return
	}

	ctx.Status(200)
}
