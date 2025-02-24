package handler

import (
	"fmt"
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

	fmt.Println(string(body))
}
