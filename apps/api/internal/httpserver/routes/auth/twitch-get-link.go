package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix/v2"
)

func (c *Auth) TwitchGetLink(ctx *gin.Context) {
	redirectUri := ctx.Query("redirectUri")

	twitchClient, err := helix.NewClient(
		&helix.Options{
			ClientID:     c.config.TwitchClientID,
			ClientSecret: c.config.TwitchSecret,
			RedirectURI:  c.config.TwitchRedirectURI,
		},
	)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	authLink := twitchClient.GetAuthorizationURL(
		&helix.AuthorizationURLParams{
			ResponseType: "code",
			Scopes:       []string{},
			ForceVerify:  true,
		},
	)

	session := c.sessionStore.GetSession(ctx.Request.Context())

	session.Set("twitchRedirectUri", redirectUri)
	session.Save()

	ctx.Redirect(302, authLink)
}
