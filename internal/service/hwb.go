package service

import (
	"com.setlog/internal/model/iata"
	"encoding/json"
	"strconv"
	"strings"

	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
)

const promptHwbVertexAI = `You are a very professional specialist in analysing inspections of consumer goods.
            Please check if the given document is a report of a consumer good inspection test.
            If it is one, then please summarize the given document and report me, if the test results in this document show a passed or failed status.
            If it is not, then please set the Status to "UNDEFINED" and the rest of the fields to null.
			Also provide me general document information, i want know the inspection report number "InspectionReportNo", the inspection date "InspectionDate" and the format of the date e.g. yyyy-mm-dd "InspectionDateFormat".
			Also give me the person who performed the inspection. This field is most often marked with "Inspected by" in the inspection report. The name of the json field is "Inspector".
            Give me a feedback as plain json (no json encasing) with fields "InspectionResult", "InspectionDate" and "InspectionReportNo"
Here is the remplate for the response json struct
{
    "InspectionResult": "PASS" | "FAIL" | "UNDEFINED",
	"InspectionReportNo": "222416610",
	"InspectionDate": "27-05-2024",
	"InspectionDateFormat": "DD.MM.YYYY",
	"Inspector": "Spark Ling"
}`

type HwbService struct {
	config *configuration.Config
	ai     *AiCommunicationService
}

func NewHwbService(config *configuration.Config) *HwbService {
	ai := NewAiCommunicationService(config)
	return &HwbService{config: config, ai: ai}
}

func (i *HwbService) AnalysePdfFile(filename string) (*model.EntityCollection, error) {
	answer, err := i.ai.GenerateContentFromPDF(filename, promptHwbVertexAI)
	result := model.HwbReportResponseVertexAi{}
	err = json.Unmarshal([]byte(answer), &result)
	if err != nil {
		return nil, err
	}
	ValidationResp := i.ConvertResponse(&result)
	return ValidationResp, nil
}

func (i *HwbService) ConvertResponse(responseVertexAI *model.HwbReportResponseVertexAi) *model.EntityCollection {
	entityCollection := model.EntityCollection{}
	context := iata.Context{Cargo: "https://onerecord.iata.org/ns/cargo#"}
	carrier := iata.Organization{
		Context: context,
		Type:    "cargo:Organization",
		Name:    responseVertexAI.CarrierName,
	}
	entityCollection.Organizations = append(entityCollection.Organizations, carrier)
	shipper := iata.Organization{
		Context: context,
		Type:    "cargo:Organization",
		Name:    responseVertexAI.ShipperName,
	}
	entityCollection.Organizations = append(entityCollection.Organizations, shipper)
	consignee := iata.Organization{
		Context: context,
		Type:    "cargo:Organization",
		Name:    responseVertexAI.ConsigneeName,
	}
	entityCollection.Organizations = append(entityCollection.Organizations, consignee)
	factoryName := "Unknown"
	if responseVertexAI.FactoryName != "" {
		factoryName = responseVertexAI.FactoryName
	}
	var itemDescriptions []string
	for _, p := range responseVertexAI.ShipmentOfPieces {
		product := iata.Product{
			Context:          context,
			Type:             "cargo:Product",
			Description:      p.ItemDescription,
			HsCode:           p.HsCode,
			UniqueIdentifier: p.ItemNumber,
			Manufacturer: iata.Manufacturer{
				Type: "cargo:Company",
				Name: factoryName,
			},
		}
		item := iata.Item{
			Context:     context,
			Type:        "cargo:Item",
			Measurement: iata.Measurement{Type: "cargo:Value", Unit: p.Unit, NumericalValue: float64(p.Quantity)},
			RawProduct:  product,
		}
		entityCollection.Products = append(entityCollection.Products, product)
		entityCollection.Items = append(entityCollection.Items, item)
		itemDescriptions = append(itemDescriptions, p.ItemDescription)
	}
	piece := iata.Piece{
		Context:              context,
		Type:                 "cargo:Piece",
		Coload:               true,
		GoodsDescription:     strings.Join(itemDescriptions, ", "),
		Upid:                 "",
		ContainedItems:       nil,
		HandlingInstructions: nil,
	}
	entityCollection.Pieces = append(entityCollection.Pieces, piece)
	totalDimensions := iata.TotalDimensions{
		Height: iata.Measurement{
			Type:           "cargo:Measurement",
			Unit:           responseVertexAI.TotalDimensions.Unit,
			NumericalValue: convertRawValueToNumericalValue(responseVertexAI.TotalDimensions.Height),
		},
		Length: iata.Measurement{
			Type:           "cargo:Measurement",
			Unit:           responseVertexAI.TotalDimensions.Unit,
			NumericalValue: convertRawValueToNumericalValue(responseVertexAI.TotalDimensions.Length)},
		Width: iata.Measurement{
			Type:           "cargo:Measurement",
			Unit:           responseVertexAI.TotalDimensions.Unit,
			NumericalValue: convertRawValueToNumericalValue(responseVertexAI.TotalDimensions.Width)},
		Volume: iata.Measurement{
			Type:           "cargo:Measurement",
			Unit:           responseVertexAI.TotalDimensions.Unit,
			NumericalValue: calculateVolume(responseVertexAI.TotalDimensions)},
	}
	shipment := iata.Shipment{
		Context:          context,
		Type:             "cargo:Shipment",
		GoodsDescription: strings.Join(itemDescriptions, ", "),
		TotalGrossWeight: iata.Measurement{Type: "cargo:Value", Unit: "KG", NumericalValue: convertRawValueToNumericalValue(responseVertexAI.TotalGrossWeight)},
		TotalDimensions:  totalDimensions,
		ShipmentOfPieces: nil,
		InvolvedParties:  nil,
	}
	entityCollection.Shipments = append(entityCollection.Shipments, shipment)
	hwb := iata.Hwb{Context: context, Type: "cargo:Waybill", WaybillNumber: responseVertexAI.Hawb, WaybillType: "house"}
	entityCollection.Hwbs = append(entityCollection.Hwbs, hwb)
	return &entityCollection
}

func convertRawValueToNumericalValue(rawValue string) float64 {
	convertedValue, _ := strconv.ParseFloat(rawValue, 64)
	return convertedValue
}

func calculateVolume(totalDimensions model.TotalDimensions) float64 {
	return convertRawValueToNumericalValue(totalDimensions.Height) * convertRawValueToNumericalValue(totalDimensions.Length) * convertRawValueToNumericalValue(totalDimensions.Width)
}
