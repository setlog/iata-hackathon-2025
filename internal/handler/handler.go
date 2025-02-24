package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"com.setlog/internal/configuration"
	"com.setlog/internal/service"
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
	token := service.NewTokenService(h.config)
	iata := service.NewIataService(h.config, token)
	fileName := req.URL.Query().Get("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("parameter 'fileName' not found.")
		fmt.Fprint(w, "parameter 'fileName' not found.")
		return
	}
	answer, err := h.hwb.AnalysePdfFile(fileName)
	iata.CreateIataData(answer)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("%v", slog.Any("error", err))
		fmt.Fprintf(w, "%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", answer)
}
