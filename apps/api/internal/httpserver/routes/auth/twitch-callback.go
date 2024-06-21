package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/nicklaw5/helix/v2"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	userrepo "github.com/satont/stream/apps/api/internal/repositories/user"
)

func (c *Auth) TwitchCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(400, gin.H{"error": "missing code"})
		return
	}

	twitchClient, err := helix.NewClientWithContext(
		ctx.Request.Context(),
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

	token, err := twitchClient.RequestUserAccessToken(code)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	twitchClient.SetUserAccessToken(token.Data.AccessToken)

	userResponse, err := twitchClient.GetUsers(&helix.UsersParams{})
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if userResponse.ErrorMessage != "" {
		ctx.JSON(500, gin.H{"error": userResponse.ErrorMessage})
		return
	}
	if len(userResponse.Data.Users) == 0 {
		ctx.JSON(500, gin.H{"error": "Cannot get authorized user"})
		return
	}

	user := userResponse.Data.Users[0]

	session := c.sessionStore.GetSession(ctx.Request.Context())

	existedUser, err := c.userRepo.FindByProviderUserID(
		ctx.Request.Context(),
		user.ID,
		userrepo.UserConnectionProviderTwitch,
	)
	if err != nil && !errors.Is(err, userrepo.ErrNotFound) {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if existedUser != nil {
		err = c.userRepo.UpdateProviderByUserID(
			ctx.Request.Context(),
			existedUser.ID,
			userrepo.UserConnectionProviderTwitch,
			userrepo.UpdateProviderByUserIdOpts{
				ProviderUserName:        user.Login,
				ProviderUserDisplayName: user.DisplayName,
				ProviderAvatar:          user.ProfileImageURL,
			},
		)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		session.Set(session_storage.USER_ID_KEY, existedUser.ID.String())
	} else {
		newUser, err := c.userRepo.Create(
			ctx.Request.Context(),
			userrepo.CreateOpts{
				Name:        user.Login,
				DisplayName: user.DisplayName,
				AvatarUrl:   user.ProfileImageURL,
				Color:       colorful.WarmColor().Hex(),
				Provider: userrepo.CreateOptsProvider{
					Provider:                userrepo.UserConnectionProviderTwitch,
					ProviderUserID:          user.ID,
					ProviderUserName:        user.Login,
					ProviderUserDisplayName: user.DisplayName,
					ProviderAvatar:          user.ProfileImageURL,
				},
			},
		)
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		session.Set(session_storage.USER_ID_KEY, newUser.ID.String())
	}

	if err := session.Save(); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	twitchRedirectUri := session.Get("twitchRedirectUri")
	if twitchRedirectUri != nil && twitchRedirectUri != "" {
		ctx.Redirect(302, twitchRedirectUri.(string))
		return
	}

	ctx.JSON(200, gin.H{session_storage.USER_ID_KEY: session.Get(session_storage.USER_ID_KEY)})
}
