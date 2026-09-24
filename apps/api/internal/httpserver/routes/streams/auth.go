package streams

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	userrepo "github.com/satont/stream/apps/api/internal/repositories/user"
)

type admissionClient struct {
	Address   string `json:"address"`
	Port      int    `json:"port"`
	RealIP    string `json:"real_ip"`
	UserAgent string `json:"user_agent"`
}

type admissionRequest struct {
	Direction string `json:"direction"`
	Protocol  string `json:"protocol"`
	Status    string `json:"status"`
	Url       string `json:"url"`
	NewUrl    string `json:"new_url"`
	Time      string `json:"time"`
}

type admissionBody struct {
	Client  admissionClient  `json:"client"`
	Request admissionRequest `json:"request"`
}

type admissionResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}

// authHandler handles AdmissionWebhooks requests from OvenMediaEngine.
// Only publishing is authorized here: the encoder url must contain a valid
// stream key (`rtmp://host/app/<channel>?key=<stream_key>`).
func (c *Streams) authHandler(ctx *gin.Context) {
	rawBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	if !c.verifySignature(rawBody, ctx.GetHeader("X-OME-Signature")) {
		c.logger.Sugar().Warnw("Invalid ome signature on admission request", "path", ctx.Request.URL.Path)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid signature"})
		return
	}

	body := admissionBody{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		c.logger.Sugar().Infow("Invalid body", "err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Streams are publicly readable, only publishing needs authorization.
	if body.Request.Status != "opening" || body.Request.Direction != "incoming" {
		ctx.JSON(http.StatusOK, admissionResponse{Allowed: true})
		return
	}

	dbUser, err := c.getPublishUser(ctx, body.Request.Url)
	if err != nil {
		c.logger.Sugar().Infow("Publish denied", "err", err, "url", body.Request.Url)
		ctx.JSON(http.StatusOK, admissionResponse{Allowed: false, Reason: err.Error()})
		return
	}

	c.logger.Sugar().Infow("Publish allowed", "channel", dbUser.Name, "user_id", dbUser.ID)
	ctx.JSON(http.StatusOK, admissionResponse{Allowed: true})
}

func (c *Streams) getPublishUser(ctx *gin.Context, rawUrl string) (*userrepo.User, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errInvalidUrl
	}

	// /app/<channel>
	pathParts := strings.Split(strings.Trim(parsedUrl.Path, "/"), "/")
	if len(pathParts) < 2 {
		return nil, errInvalidUrl
	}
	channelName := pathParts[len(pathParts)-1]

	streamKey, err := uuid.Parse(parsedUrl.Query().Get("key"))
	if err != nil {
		return nil, errInvalidStreamKey
	}

	dbUser, err := c.userRepo.FindByStreamKey(ctx.Request.Context(), streamKey)
	if err != nil {
		return nil, errInvalidStreamKey
	}

	if dbUser.Banned {
		return nil, errBanned
	}

	if dbUser.Name != channelName {
		return nil, errChannelMismatch
	}

	return dbUser, nil
}

func (c *Streams) verifySignature(body []byte, signature string) bool {
	secret := c.config.OmeAdmissionSecret
	if secret == "" {
		return true
	}

	if signature == "" {
		return false
	}

	actual, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		actual, err = base64.RawURLEncoding.DecodeString(signature)
		if err != nil {
			return false
		}
	}

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(body)

	return hmac.Equal(actual, mac.Sum(nil))
}
