package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadAWS(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-southeast-1"),
		config.WithSharedConfigProfile(os.Getenv("AWS_PROFILE")),
	)
}