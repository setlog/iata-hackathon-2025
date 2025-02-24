package handler

import (
	"com.setlog/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func (h *AiHandler) Json2Iata(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("error while reading body", slog.Any("error", err))
		return
	}
	defer req.Body.Close()

	conv, err := h.ConvertToResponseVertexAi(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("error while converting to response", slog.Any("error", err))
		return
	}

	if conv.IsHawb {
		h.hwb.ConvertResponse(conv)
	}
}

func (h *AiHandler) Json2IataAll(writer http.ResponseWriter, request *http.Request) {

	files, err := os.ReadDir("./ai-output")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		slog.Error("error while reading directory", slog.Any("error", err))
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := fmt.Sprintf("./ai-output/%s", file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				slog.Error("error while reading file", slog.String("file", file.Name()), slog.Any("error", err))
				continue
			}

			response, err := h.ConvertToResponseVertexAi(content)
			if err != nil {
				slog.Error("error while converting to response", slog.String("file", file.Name()), slog.Any("error", err))
				continue
			}
			if response.IsHawb {
				h.hwb.ConvertResponse(response)
			}
		}
	}
}

func (h *AiHandler) ConvertToResponseVertexAi(payload []byte) (*model.HwbReportResponseVertexAi, error) {
	var resp *model.HwbReportResponseVertexAi
	err := json.Unmarshal(payload, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
