package streams

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/grafov/m3u8"
	"github.com/imroc/req/v3"
	"golang.org/x/sync/errgroup"
)

func (c *Streams) indexHandler(ctx *gin.Context) {
	refererUrl, err := url.Parse(ctx.Request.Referer())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referer"})
		return
	}

	channelID := ctx.Param("channelID")
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	user, err := c.getUserById(ctx.Request.Context(), parsedChannelID)
	if err != nil {
		c.logger.Sugar().Errorw(
			"Failed to get user by ID",
			"err",
			err,
			"channel_id",
			parsedChannelID,
		)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
		return
	}

	playlist, err := c.buildPlaylist(
		ctx.Request.Context(),
		refererUrl,
		user.StreamKey.String(),
		user.ID.String(),
	)
	if err != nil {
		c.logger.Sugar().Errorw(
			"Failed to build playlist",
			"err",
			err,
			"user_id",
			user.ID,
			"referer",
			ctx.Request.Referer(),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build playlist"})
		return
	}

	ctx.String(http.StatusOK, playlist)
}

func (c *Streams) buildPlaylist(
	ctx context.Context,
	refererUrl *url.URL,
	streamKey,
	userId string,
) (
	string,
	error,
) {
	baseUrl := fmt.Sprintf(
		"%s://%s/api/streams/read",
		refererUrl.Scheme,
		refererUrl.Host,
	)

	// first quality is source
	var resolutions []string
	if c.config.FeatureDisableQuality {
		resolutions = []string{""}
	} else {
		resolutions = []string{"", "720p_", "480p_", "360p_"}
	}

	masterPlaylist := m3u8.NewMasterPlaylist()
	masterPlaylist.SetIndependentSegments(true)
	masterPlaylist.SetVersion(6)
	var plmu sync.Mutex
	var errwg errgroup.Group
	for _, resolution := range resolutions {
		resolution := resolution
		uri := fmt.Sprintf(
			"%s/%s%s/stream.m3u8",
			baseUrl,
			resolution,
			userId,
		)

		errwg.Go(
			func() error {
				streamReqUri := fmt.Sprintf(
					"%s/%s%s/index.m3u8",
					c.config.MediaMtxAddr+":8888",
					resolution,
					streamKey,
				)

				resp, err := req.SetContext(ctx).Get(
					fmt.Sprintf(
						"%s/%s%s/index.m3u8",
						c.config.MediaMtxAddr+":8888",
						resolution,
						streamKey,
					),
				)
				if err != nil {
					return fmt.Errorf("failed to fetch m3u8 %s: %w", streamReqUri, err)
				}
				p, _, err := m3u8.DecodeFrom(resp.Body, true)
				if err != nil {
					return fmt.Errorf("failed to decode m3u8 %s: %w", streamReqUri, err)
				}

				masterpl := p.(*m3u8.MasterPlaylist)

				for _, variant := range masterpl.Variants {
					variant.URI = uri

					name := strings.ReplaceAll(variant.VariantParams.Resolution, "_", "")

					plmu.Lock()
					variant.Video = fmt.Sprintf("%s%s", name, variant.Video)
					variant.VariantParams.Video = name
					if name == "" {
						variant.VariantParams.Video = "source"
					}
					masterPlaylist.Variants = append(masterPlaylist.Variants, variant)
					plmu.Unlock()
				}

				return nil
			},
		)
	}

	if err := errwg.Wait(); err != nil {
		return "", err
	}

	return masterPlaylist.String(), nil
}
