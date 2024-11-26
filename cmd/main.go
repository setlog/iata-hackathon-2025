package main

import (
	conf "com.setlog/internal/configuration"
	"com.setlog/internal/handler"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	err, configuration := conf.NewConfig()
	if err != nil {
		slog.Error("Error parsing .env file", err)
		return
	}
	hand := handler.NewAiHandler(configuration)

	r := mux.NewRouter()
	//r.Use(middleware.AuthMiddleware(configuration))
	r.HandleFunc("/producttestreportanalysis", hand.ProductTestHandlerFunc)
	r.HandleFunc("/inspectionreportanalysis", hand.InspectionHandlerFunc)

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
