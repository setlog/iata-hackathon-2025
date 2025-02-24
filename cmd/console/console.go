package main

import (
	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
	"com.setlog/internal/service"
	"encoding/json"
	"fmt"
	//"github.com/spf13/viper"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
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
	statistic := model.Statistic{}
	var payloads []FilePayload
	err = filepath.Walk("../example/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		payload, err := os.ReadFile(path)
		payloads = append(payloads, FilePayload{Payload: payload})
		return err
	})
	//payload, err := os.ReadFile("../example/ai_reponse.json")
	if err != nil {
		return
	}
	for _, payload := range payloads {
		var resp *model.HwbReportResponseVertexAi
		err = json.Unmarshal(payload.Payload, &resp)
		if err != nil {
			return
		}

		if isPayloadValid(&statistic, resp) {
			conv := hw.ConvertResponse(resp)

			err = iata.CreateIataData(conv)
			if err != nil {
				slog.Error(err.Error())
				return
			}
		}

	}
	fmt.Printf("%#v", statistic)

}

func isPayloadValid(statistic *model.Statistic, resp *model.HwbReportResponseVertexAi) bool {
	pieceFailed := false
	grossWeightMissing := false

	statistic.TotalFiles++
	for _, piece := range resp.ShipmentOfPieces {
		itemNumberMissing := false
		quantityMissing := false
		manufacturerMissing := false
		statistic.TotalNumberOfItems++
		if piece.ItemNumber == "" {
			statistic.ItemNumberMissing++
			itemNumberMissing = true
		}
		if piece.Quantity == 0 {
			statistic.QuantityMissing++
			quantityMissing = true
		}
		if piece.Manufacturer == "" {
			statistic.ManufacturerMissing++
			manufacturerMissing = true
		}
		if itemNumberMissing || quantityMissing || manufacturerMissing {
			statistic.NumberOfFailedItems++
			pieceFailed = true
		}
	}

	if resp.TotalGrossWeight == "" {
		grossWeightMissing = true
	}

	if pieceFailed || grossWeightMissing {
		statistic.FailedFiles++
		return false
	}
	return true
}

type FilePayload struct {
	Payload []byte
}
