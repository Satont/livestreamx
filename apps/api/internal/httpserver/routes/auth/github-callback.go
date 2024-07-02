package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	session_storage "github.com/satont/stream/apps/api/internal/httpserver/session-storage"
	userrepo "github.com/satont/stream/apps/api/internal/repositories/user"
)

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (c *Auth) GithubCallback(ctx *gin.Context) {
	session := c.sessionStore.GetSession(ctx.Request.Context())
	redirectUri := session.Get("redirectUri")
	if redirectUri == nil {
		ctx.JSON(400, gin.H{"error": "missing redirectUri"})
		return
	}

	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(400, gin.H{"error": "missing code"})
		return
	}

	githubTokenData := &githubTokenResponse{}
	tokenResp, err := req.
		SetBody(
			map[string]any{
				"client_id":     c.config.GithubClientID,
				"client_secret": c.config.GithubClientSecret,
				"code":          code,
			},
		).
		SetHeader("Accept", "application/json").
		SetSuccessResult(githubTokenData).
		Post("https://github.com/login/oauth/access_token")
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !tokenResp.IsSuccessState() {
		ctx.JSON(
			500,
			gin.H{"error": fmt.Sprintf("github token request failed: %s", tokenResp.String())},
		)
		return
	}

	githubUser := &githubUserResponse{}
	userResp, err := req.
		SetHeader("Authorization", "token "+githubTokenData.AccessToken).
		SetHeader("Accept", "application/json").
		SetSuccessResult(githubUser).
		Get("https://api.github.com/user")
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if !userResp.IsSuccessState() {
		ctx.JSON(
			500,
			gin.H{"error": fmt.Sprintf("github user request failed: %s", userResp.String())},
		)
		return
	}

	dbUser, err := c.callback(
		ctx.Request.Context(),
		callbackOpts{
			Provider:                userrepo.UserConnectionProviderGithub,
			ProviderUserID:          strconv.Itoa(githubUser.Id),
			ProviderUserName:        githubUser.Login,
			ProviderUserDisplayName: githubUser.Login,
			ProviderAvatar:          githubUser.AvatarUrl,
			Email:                   &githubUser.Email,
		},
	)
	if err != nil {
		c.logger.Sugar().Error(err)
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session.Set(session_storage.USER_ID_KEY, dbUser.ID.String())
	if err := session.Save(); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if redirectUri != nil && redirectUri != "" {
		ctx.Redirect(302, redirectUri.(string))
		return
	}

	ctx.JSON(200, gin.H{"token": githubTokenData.AccessToken[0:3]})
}

type githubUserResponse struct {
	Login                   string    `json:"login"`
	Id                      int       `json:"id"`
	NodeId                  string    `json:"node_id"`
	AvatarUrl               string    `json:"avatar_url"`
	GravatarId              string    `json:"gravatar_id"`
	Url                     string    `json:"url"`
	HtmlUrl                 string    `json:"html_url"`
	FollowersUrl            string    `json:"followers_url"`
	FollowingUrl            string    `json:"following_url"`
	GistsUrl                string    `json:"gists_url"`
	StarredUrl              string    `json:"starred_url"`
	SubscriptionsUrl        string    `json:"subscriptions_url"`
	OrganizationsUrl        string    `json:"organizations_url"`
	ReposUrl                string    `json:"repos_url"`
	EventsUrl               string    `json:"events_url"`
	ReceivedEventsUrl       string    `json:"received_events_url"`
	Type                    string    `json:"type"`
	SiteAdmin               bool      `json:"site_admin"`
	Name                    string    `json:"name"`
	Company                 string    `json:"company"`
	Blog                    string    `json:"blog"`
	Location                string    `json:"location"`
	Email                   string    `json:"email"`
	Hireable                bool      `json:"hireable"`
	Bio                     string    `json:"bio"`
	TwitterUsername         string    `json:"twitter_username"`
	PublicRepos             int       `json:"public_repos"`
	PublicGists             int       `json:"public_gists"`
	Followers               int       `json:"followers"`
	Following               int       `json:"following"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	PrivateGists            int       `json:"private_gists"`
	TotalPrivateRepos       int       `json:"total_private_repos"`
	OwnedPrivateRepos       int       `json:"owned_private_repos"`
	DiskUsage               int       `json:"disk_usage"`
	Collaborators           int       `json:"collaborators"`
	TwoFactorAuthentication bool      `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		PrivateRepos  int    `json:"private_repos"`
		Collaborators int    `json:"collaborators"`
	} `json:"plan"`
}
