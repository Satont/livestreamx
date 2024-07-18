package streams

import (
	"net/url"
	"regexp"

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

var (
	authKeyRegexp            = regexp.MustCompile("^(?P<params>.+)$")
	authKeyRegexpWithQuality = regexp.MustCompile("^(?P<quality>.+p_)(?P<params>.+)$")
)

func (c *Streams) authHandler(ctx *gin.Context) {
	body := authReq{}
	if err := ctx.BindJSON(&body); err != nil {
		c.logger.Sugar().Infow("Invalid body", "err", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if body.Action == "read" {
		ctx.Status(200)
		return
	}

	if body.Action == "publish" {
		query, err := url.ParseQuery(body.Query)
		if err != nil {
			c.logger.Sugar().Infow("Invalid query", "err", err)
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		streamKey, err := uuid.Parse(query.Get("key"))
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
