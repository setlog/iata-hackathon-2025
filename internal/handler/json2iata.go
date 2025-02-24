package handler

import (
	"com.setlog/internal/model"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func (h *AiHandler) Json2Iata(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("error while reading body", slog.Any("error", err))
		return
	}
	defer req.Body.Close()

	var resp *model.HwbReportResponseVertexAi
	err = json.Unmarshal(body, &resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("error while reading body", slog.Any("error", err))
		return
	}
}
