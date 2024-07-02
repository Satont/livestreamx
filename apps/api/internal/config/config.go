package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv string `envconfig:"APP_ENV" default:"development"`

	RedisURL    string `envconfig:"REDIS_URL"`
	PostgresURL string `envconfig:"POSTGRES_URL"`
	NatsURL     string `envconfig:"NATS_URL" default:"nats://localhost:8223" required:"true"`

	ApiPort          int    `required:"false"  envconfig:"API_PORT" default:"1337"`
	ApiSessionSecret string `envconfig:"API_SESSION_SECRET"`
	MediaMtxAddr     string `envconfig:"MEDIA_MTX_ADDR" required:"true"`
	ThumbnailsURI    string `envconfig:"THUMBNAILS_URI" required:"true"`

	S3Host        string `envconfig:"S3_HOST"`
	S3AccessToken string `envconfig:"S3_ACCESS_TOKEN"`
	S3SecretToken string `envconfig:"S3_SECRET_TOKEN"`
	S3Region      string `envconfig:"S3_REGION"`
	S3Bucket      string `envconfig:"S3_BUCKET"`

	TwitchClientID    string `envconfig:"TWITCH_CLIENT_ID"`
	TwitchSecret      string `envconfig:"TWITCH_SECRET"`
	TwitchRedirectURI string `envconfig:"TWITCH_REDIRECT_URI"`

	GithubClientID     string `required:"false" envconfig:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `required:"false" envconfig:"GITHUB_CLIENT_SECRET"`
	GithubRedirectURI  string `required:"false" envconfig:"GITHUB_REDIRECT_URI"`
}

func New() (Config, error) {
	cfg := Config{}

	wd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	envPath := filepath.Join(wd, ".env")

	_ = godotenv.Load(envPath)

	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
