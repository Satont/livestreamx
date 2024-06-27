package auth

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func (c *Auth) GithubGetLink(ctx *gin.Context) {
	redirectUri := ctx.Query("redirectUri")

	u := url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/login/oauth/authorize",
	}

	q := u.Query()
	q.Set("client_id", c.config.GithubClientID)
	q.Set("redirect_uri", c.config.GithubRedirectURI)
	q.Set("scope", "read:user user:email")

	u.RawQuery = q.Encode()

	authLink := u.String()

	session := c.sessionStore.GetSession(ctx.Request.Context())

	session.Set("redirectUri", redirectUri)
	session.Save()

	ctx.Redirect(302, authLink)
}
