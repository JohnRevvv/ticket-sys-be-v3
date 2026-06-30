package controllers

import (
	"context"
	"fmt"
	"time"

	"ideyanale-be/pkg/config"
	"ideyanale-be/pkg/services/upload/model"
	"ideyanale-be/pkg/services/upload/script"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v3"
)

type UploadController struct{}

func (uc *UploadController) Upload(c fiber.Ctx) error {
	ctx := context.Background()

	// Get file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "file is required",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to open file",
		})
	}
	defer file.Close()

	// AWS config
	cfg, err := config.LoadAWS(ctx)
	if err != nil {
		return c.Status(500).JSON(err)
	}

	client := s3.NewFromConfig(cfg)
	s3Script := scripts.S3Script{Client: client}

	bucket := "idiyanale-bucket"
	key := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), fileHeader.Filename)

	// Upload
	if err := s3Script.Upload(ctx, bucket, key, file); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "upload failed",
			"error":   err.Error(),
		})
	}

	resp := models.UploadResponse{
		FileName: fileHeader.Filename,
		Bucket:   bucket,
		Key:      key,
		URL:      fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key),
	}

	return c.JSON(resp)
}
