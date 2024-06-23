package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv             string `envconfig:"APP_ENV" default:"development"`
	RedisURL           string `envconfig:"REDIS_URL"`
	PostgresURL        string `envconfig:"POSTGRES_URL"`
	TwitchClientID     string `envconfig:"TWITCH_CLIENT_ID"`
	TwitchSecret       string `envconfig:"TWITCH_SECRET"`
	TwitchRedirectURI  string `envconfig:"TWITCH_REDIRECT_URI"`
	ApiPort            int    `required:"false"  envconfig:"API_PORT" default:"1337"`
	ApiSessionSecret   string `envconfig:"API_SESSION_SECRET"`
	SevenTVEmoteSetID  string `envconfig:"SEVEN_TV_EMOTE_SET_ID"`
	StreamPathInfoAddr string `envconfig:"STREAM_PATH_INFO_ADDR"`

	S3Host        string `envconfig:"S3_HOST"`
	S3AccessToken string `envconfig:"S3_ACCESS_TOKEN"`
	S3SecretToken string `envconfig:"S3_SECRET_TOKEN"`
	S3Region      string `envconfig:"S3_REGION"`
	S3Bucket      string `envconfig:"S3_BUCKET"`
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
