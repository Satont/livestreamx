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
	channelID := ctx.Param("channelID")
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}
	refererUrl, err := url.Parse(ctx.Request.Referer())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referer"})
		return
	}

	user, err := c.getUserById(ctx.Request.Context(), parsedChannelID)
	if err != nil {
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
		c.logger.Sugar().Errorf("Failed to build playlist: %v", err)
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

	// first is source
	resolutions := []string{"", "720p_", "480p_", "360p_"}
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
				resp, err := req.SetContext(ctx).Get(
					fmt.Sprintf(
						"%s/%s%s/index.m3u8",
						c.config.MediaMtxAddr+":8888",
						resolution,
						streamKey,
					),
				)
				if err != nil {
					return err
				}
				p, _, err := m3u8.DecodeFrom(resp.Body, true)
				if err != nil {
					return err
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
					// variant.Alternatives = append(
					// 	[]*m3u8.Alternative{}, &m3u8.Alternative{
					// 		GroupId:    name,
					// 		URI:        uri,
					// 		Type:       "VIDEO",
					// 		Name:       name,
					// 		Default:    true,
					// 		Autoselect: "YES",
					// 	},
					// )
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
