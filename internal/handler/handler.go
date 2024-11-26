package handler

import (
	"com.setlog/internal/configuration"
	"com.setlog/internal/service"
	"fmt"
	"log/slog"
	"net/http"
)

type AiHandler struct {
	product    *service.ProductTestService
	inspection *service.InspectionService
	config     *configuration.Config
}

func NewAiHandler(config *configuration.Config) *AiHandler {
	product := service.NewProductTestService(config)
	inspection := service.NewInspectionService(config)
	return &AiHandler{product: product, inspection: inspection, config: config}
}

func (h *AiHandler) ProductTestHandlerFunc(w http.ResponseWriter, req *http.Request) {
	fileName := req.URL.Query().Get("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("parameter 'fileName' not found.")
		fmt.Fprint(w, "parameter 'fileName' not found.")
		return
	}
	answer, err := h.inspection.AnalysePdfFile(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("%v", slog.Any("error", err))
		fmt.Fprintf(w, "%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(answer))
}

func (h *AiHandler) InspectionHandlerFunc(w http.ResponseWriter, req *http.Request) {
	fileName := req.URL.Query().Get("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("parameter 'fileName' not found.")
		fmt.Fprint(w, "parameter 'fileName' not found.")
		return
	}

	answer, err := h.inspection.AnalysePdfFile(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("%v", slog.Any("error", err))
		fmt.Fprintf(w, "%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(answer))
}
