package otelsdk

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type S3Client struct {
	Client *s3.Client
}

func NewS3Client(ctx context.Context) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)
	s3Client := s3.NewFromConfig(cfg)

	return &S3Client{Client: s3Client}, nil
}
