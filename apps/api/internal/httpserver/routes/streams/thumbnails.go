package streams

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (c *Streams) thumbnailsHandler(ctx *gin.Context) {
	channelID := ctx.Param("channelID")
	channelID = strings.ReplaceAll(channelID, "/", "")

	channelUuid, err := uuid.Parse(channelID)
	if err != nil {
		c.logger.Sugar().Error(err)
		ctx.String(400, "invalid channel id")
		return
	}

	dbChannel, err := c.userRepo.FindByID(ctx.Request.Context(), channelUuid)
	if err != nil {
		c.logger.Sugar().Error(err)
		ctx.String(404, "channel not found")
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		c.logger.Sugar().Error(err)
		ctx.String(500, "internal error happened on getting filepath for thumbnail")
		return
	}

	var thumbnailsRoot string
	if c.config.AppEnv != "production" {
		thumbnailsRoot = filepath.Join(wd, "thumbnails")
	} else {
		thumbnailsRoot = "/thumbnails"
	}

	thumbnailPath := filepath.Join(thumbnailsRoot, dbChannel.StreamKey.String()+".jpg")
	c.logger.Sugar().Info("thumbnailPath: ", thumbnailPath)

	ctx.File(thumbnailPath)
}
