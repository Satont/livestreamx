package s3

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/satont/stream/apps/api/internal/config"
	"go.uber.org/fx"
)

func New(cfg config.Config, lc fx.Lifecycle) (*minio.Client, error) {
	var creds *credentials.Credentials
	if cfg.AppEnv != "production" {
		creds = credentials.NewStaticV4("minio", "minio-password", "")
	} else {
		creds = credentials.NewStaticV4(cfg.S3AccessToken, cfg.S3SecretToken, "")
	}

	client, err := minio.New(
		cfg.S3Host,
		&minio.Options{
			Creds:  creds,
			Region: cfg.S3Region,
			Secure: cfg.AppEnv == "production",
		},
	)
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				buckets, err := client.ListBuckets(ctx)
				if err != nil {
					return fmt.Errorf("cannot list buckets: %w", err)
				}

				bucketExists := false
				for _, bucket := range buckets {
					if bucket.Name == cfg.S3Bucket {
						bucketExists = true
						break
					}
				}

				if !bucketExists {
					err = client.MakeBucket(ctx, cfg.S3Bucket, minio.MakeBucketOptions{})
					if err != nil {
						return fmt.Errorf("cannot create bucket: %w", err)
					}
				}

				// we use cloudflare r2, which doesnt support this operation
				if cfg.AppEnv != "production" {
					err = client.SetBucketPolicy(
						ctx,
						cfg.S3Bucket,
						`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": ["*"]
				},
				"Action": ["s3:GetObject"],
				"Resource": [
					"arn:aws:s3:::`+cfg.S3Bucket+`/**"
				]
			}
		]
	}`,
					)

					if err != nil {
						return fmt.Errorf("cannot set bucket policy: %w", err)
					}
				}

				return nil
			},
			OnStop: nil,
		},
	)

	return client, nil
}
