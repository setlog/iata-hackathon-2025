package handler

import (
	"cloud.google.com/go/storage"
	"com.setlog/internal/configuration"
	"com.setlog/internal/service"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type AiHandler struct {
	hwb    *service.HwbService
	config *configuration.Config
}

func NewAiHandler(config *configuration.Config) *AiHandler {
	hwb := service.NewHwbService(config)
	return &AiHandler{hwb: hwb, config: config}
}

func (h *AiHandler) HwbReportHandlerFunc(w http.ResponseWriter, req *http.Request) {
	fileName := req.URL.Query().Get("fileName")
	isOk := h.verifyFileName(fileName)
	if !isOk {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "parameter 'fileName' not found.")
		return
	}

	err := h.parseFileByFileName(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("error while parsing file", slog.Any("error", err))
		_, _ = fmt.Fprint(w, "error while parsing file")
		return
	}
}

func (h *AiHandler) HwbReportHandlerFuncAll(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("could create storage client", err)
		return
	}
	defer client.Close()

	bucketName := strings.TrimPrefix(os.Getenv("GCLOUD_BUCKETNAME"), "gs://")
	bucket := client.Bucket(bucketName)

	it := bucket.Objects(ctx, nil)

	fmt.Println("Iterating over objects in bucket...")
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			fmt.Println("couldnt iterate objects", err)
			return
		}

		fileName := attrs.Name
		isOk := h.verifyFileName(fileName)
		if !isOk {
			continue
		}

		err = h.parseFileByFileName(fileName)
		if err != nil {
			fmt.Println("error while parsing file", err)
			continue
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("done")
}

func (h *AiHandler) verifyFileName(fileName string) bool {
	isOk := fileName != ""
	if !isOk {
		slog.Warn("parameter 'fileName' not found.")
	}
	return isOk
}

func (h *AiHandler) parseFileByFileName(fileName string) error {
	return h.hwb.AnalysePdfFile(fileName)
}
