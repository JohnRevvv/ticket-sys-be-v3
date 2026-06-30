package scripts

import (
	"context"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Script struct {
	Client *s3.Client
}

func (s *S3Script) Upload(ctx context.Context, bucket, key string, file multipart.File) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})
	return err
}