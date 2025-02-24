package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	conf "com.setlog/internal/configuration"
	"com.setlog/internal/handler"
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	err, configuration := conf.NewConfig()
	if err != nil {
		slog.Error("Error parsing .env file", err)
		return
	}
	hand := handler.NewAiHandler(configuration)

	r := mux.NewRouter()
	r.HandleFunc("/hwbreportanalysis/all", hand.HwbReportHandlerFuncAll)
	r.HandleFunc("/hwbreportanalysis", hand.HwbReportHandlerFunc)

	port := 8080
	slog.Info("Service has started", slog.Int("port", port))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
