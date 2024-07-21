package streams

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (c *Streams) thumbnailsHandler(ctx *gin.Context) {
	channelName := ctx.Param("channelName")
	dbChannel, err := c.userRepo.FindByName(ctx.Request.Context(), channelName)

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

	thumbnailPath := filepath.Join(thumbnailsRoot, dbChannel.Name+".jpg")

	ctx.File(thumbnailPath)
}
