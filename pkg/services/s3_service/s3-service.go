package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	Client *s3.Client
	Bucket string
}

const MaxFileSize int64 = 10 * 1024 * 1024 // 10 MB

var allowedExtensions = map[string]bool{
	".pdf":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".txt":  true,
}

var allowedMimeTypes = map[string]bool{
	"application/pdf": true,

	"image/png":  true,
	"image/jpeg": true,

	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,

	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,

	"text/plain": true,
}

func (s *S3Service) Upload(fileHeader *multipart.FileHeader, ticketID string) (string, string, error) {

	// =========================================
	// Validate file size
	// =========================================

	if fileHeader.Size > MaxFileSize {
		return "", "", fmt.Errorf(
			"%s exceeds the maximum allowed size of 10 MB",
			fileHeader.Filename,
		)
	}

	// =========================================
	// Validate extension
	// =========================================

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	if !allowedExtensions[ext] {
		return "", "", fmt.Errorf(
			"%s has an unsupported file extension",
			fileHeader.Filename,
		)
	}

	// =========================================
	// Open file
	// =========================================

	file, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// =========================================
	// Detect actual MIME type
	// =========================================

	buffer := make([]byte, 512)

	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", "", err
	}

	contentType := http.DetectContentType(buffer[:n])

	if !allowedMimeTypes[contentType] {
		return "", "", fmt.Errorf(
			"%s has an invalid file type (%s)",
			fileHeader.Filename,
			contentType,
		)
	}

	// Reset file pointer before upload
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", "", err
	}

	// =========================================
	// Upload to S3
	// =========================================

	fileName := sanitizeFileName(fileHeader.Filename)

	objectKey := fmt.Sprintf(
		"idiyanale/attachments/%s/%d_%s",
		ticketID,
		time.Now().UnixNano(),
		fileName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", "", err
	}

	return fileName, objectKey, nil
}

func NewS3Service() (*S3Service, error) {
	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, err
	}

	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("AWS_S3_BUCKET is not configured")
	}

	return &S3Service{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucket,
	}, nil
}

// Upload uploads a file to S3.
// Returns:
//   - original filename
//   - S3 object key
// func (s *S3Service) Upload(fileHeader *multipart.FileHeader, ticketID string) (string, string, error) {

// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		return "", "", err
// 	}
// 	defer file.Close()

// 	fileName := sanitizeFileName(fileHeader.Filename)

// 	objectKey := fmt.Sprintf(
// 		"idiyanale/attachments/%s/%d_%s",
// 		ticketID,
// 		time.Now().UnixNano(),
// 		fileName,
// 	)

// 	contentType := fileHeader.Header.Get("Content-Type")
// 	if contentType == "" {
// 		contentType = "application/octet-stream"
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	_, err = s.Client.PutObject(ctx, &s3.PutObjectInput{
// 		Bucket:      aws.String(s.Bucket),
// 		Key:         aws.String(objectKey),
// 		Body:        file,
// 		ContentType: aws.String(contentType),
// 	})

// 	if err != nil {
// 		return "", "", err
// 	}

// 	return fileName, objectKey, nil
// }

func (s *S3Service) GeneratePresignedDownloadURL(key string, expiration time.Duration) (string, error) {

	presignClient := s3.NewPresignClient(s.Client)

	resp, err := presignClient.PresignGetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(expiration),
	)

	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

func (s *S3Service) Delete(key string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	return err
}

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}