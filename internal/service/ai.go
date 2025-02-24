package service

import (
	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
	"com.setlog/internal/configuration"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

type AiCommunicationService struct {
	config *configuration.Config
}

func NewAiCommunicationService(config *configuration.Config) *AiCommunicationService {
	return &AiCommunicationService{config: config}
}

func (a *AiCommunicationService) GenerateContentFromPDF(fileName string, prompt string) (string, error) {
	if err := a.waitUntilFileExists(a.config.GcBucketName, fileName); err != nil {
		return "", err
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, a.config.GcProjectId, a.config.GcLocation)
	if err != nil {
		return "", fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(a.config.AiModel)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  fmt.Sprintf("%s/%s", a.config.GcBucketName, fileName),
	}

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	res, err := model.GenerateContent(ctx, part, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from model")
	}

	return fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0]), nil
}

func (a *AiCommunicationService) waitUntilFileExists(bucketName, fileName string) error {
	bucketName, _ = strings.CutPrefix(bucketName, "gs://")
	parts := strings.SplitN(bucketName, "/", 2)
	if len(parts) > 1 {
		bucketName = parts[0]
		fileName = parts[1] + "/" + fileName
	}
	timeOutInSeconds := 60
	for i := 0; i < timeOutInSeconds; i++ {
		bExists, err := a.isFileExist(bucketName, fileName)
		if bExists {
			return nil
		}
		if err != nil {
			slog.Warn("cannot read the cloud file", slog.String("bucket", bucketName), slog.String("file", fileName), slog.Any("error", err))
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout occured while reading the bucket %s, file %s", bucketName, fileName)
}

func (a *AiCommunicationService) isFileExist(bucketName, fileName string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Get the object handle
	obj := client.Bucket(bucketName).Object(fileName)

	// Try to get the attributes of the object
	_, err = obj.Attrs(ctx)
	if err != nil {
		// If the error is a "NotFound" error, the file doesn't exist
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		// Other errors should be handled accordingly
		return false, fmt.Errorf("failed to get object attributes: %v", err)
	}

	// The object exists
	return true, nil
}
