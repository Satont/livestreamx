package streams

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/satont/stream/apps/api/internal/repositories/user"
)

var (
	proxyChannelRegexp            = regexp.MustCompile("/(?P<channel_id>.+)/(?P<params>.+)")
	proxyChannelRegexpWithQuality = regexp.MustCompile("/(?P<quality>.+p_)(?P<channel_id>.+)/(?P<params>.+)")
	proxyChannelRegexpNames       = proxyChannelRegexp.SubexpNames()
)

func (c *Streams) reverseProxy(target string) gin.HandlerFunc {
	remote, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)

	return func(ctx *gin.Context) {
		p := ctx.Param("regex")
		regexMatch := proxyChannelRegexp.FindStringSubmatch(p)
		regexMatchWithQuality := proxyChannelRegexpWithQuality.FindStringSubmatch(p)
		if len(regexMatch) < 2 && len(regexMatchWithQuality) < 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// refererUrl, err := url.Parse(ctx.Request.Referer())
		// if err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referer"})
		// 	return
		// }

		var channelId, params, quality string
		if len(regexMatchWithQuality) > 1 {
			quality, channelId, params = regexMatchWithQuality[1], regexMatchWithQuality[2], regexMatchWithQuality[3]
		} else {
			channelId, params = regexMatch[1], regexMatch[2]
		}

		parsedChannelId, err := uuid.Parse(channelId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
			return
		}

		user, err := c.getUserById(ctx.Request.Context(), parsedChannelId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
			return
		}

		// if params == "index.m3u8" && len(regexMatchWithQuality) == 0 {
		// 	playlist, err := c.buildPlaylist(
		// 		ctx.Request.Context(),
		// 		refererUrl,
		// 		user.StreamKey.String(),
		// 		user.ID.String(),
		// 	)
		// 	if err != nil {
		// 		c.logger.Sugar().Errorf("Failed to build playlist: %v", err)
		// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build playlist"})
		// 		return
		// 	}
		//
		// 	ctx.String(http.StatusOK, playlist)
		// 	return
		// }

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = fmt.Sprintf("/%s%s/%s", quality, user.StreamKey, params)
			req.Host = remote.Host
			req.URL.RawPath = fmt.Sprintf("/%s%s/%s", quality, user.StreamKey, params)
		}

		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

func (c *Streams) getUserById(ctx context.Context, channelName uuid.UUID) (*user.User, error) {
	cacheKey := "streams:cache:channel-by-id:" + channelName.String()

	cachedBytes, err := c.redis.Get(ctx, cacheKey).Bytes()
	if len(cachedBytes) > 0 {
		user := &user.User{}
		if err := json.Unmarshal(cachedBytes, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	user, err := c.userRepo.FindByID(ctx, channelName)
	if err != nil {
		return nil, err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		cacheKey,
		userBytes,
		5*time.Minute,
	).Err(); err != nil {
		return nil, err
	}

	return user, nil
}
