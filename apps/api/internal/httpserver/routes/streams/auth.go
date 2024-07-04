package streams

import (
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
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if body.Action == "read" {
		ctx.Status(200)
		return
	}

	if body.Action == "publish" {
		var key string
		keyMatch := authKeyRegexp.FindStringSubmatch(body.Path)
		keyMatchWithQuality := authKeyRegexpWithQuality.FindStringSubmatch(body.Path)
		if len(keyMatch) < 2 && len(keyMatchWithQuality) < 2 {
			ctx.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if len(keyMatchWithQuality) > 1 {
			key = keyMatchWithQuality[2]
		} else {
			key = keyMatch[1]
		}

		streamKey, err := uuid.Parse(key)
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
