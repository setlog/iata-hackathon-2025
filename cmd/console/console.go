package main

import (
	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
	"com.setlog/internal/service"
	"encoding/json"
	"log/slog"
	"os"
)

func main() {

	err, conf := configuration.NewConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	token := service.NewTokenService(conf)
	iata := service.NewIataService(conf, token)
	hw := service.NewHwbService(conf)
	payload, err := os.ReadFile("./example/ai_reponse.json")
	if err != nil {
		return
	}
	var resp *model.HwbReportResponseVertexAi
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		return
	}

	conv := hw.ConvertResponse(resp)

	iata.CreateIataData(conv)
	if err != nil {
		slog.Error(err.Error())
		return
	}

}
