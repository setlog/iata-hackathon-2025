package service

import (
	"com.setlog/internal/model/iata"
	"fmt"
	"os"
	"strconv"
	"strings"

	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
)

const promptHwbVertexAI = `You are an expert air freight forwarder. Your task is to generate a HAWB (House Air Waybill) in JSON format from the provided PDF file content.

**Instructions:**

1. **Parse PDF Content:** Extract the relevant information from the provided PDF content.
2. **Required Fields:** Ensure the following required fields are present: "hawb", "totalGrossWeight", "volume", and "cargoName".If any of these fields are missing, generate an error message as described below.
3. **Optional Fields:** For optional fields that are not found in the PDF, leave the corresponding JSON values as empty strings or null.
4. **Data Format Variations:** Handle variations in data formats as described below.
5. **Handling of Images: ** Pay close attention to any embedded images, specifically JPEGs, and use OCR to accurately capture any text within them. Return the results as a JSON object.
6. **JSON Structure:**  The JSON output must strictly adhere to the following structure: json{
   "isHawb": <decide wether the document is a HAWB or not, true or false>"
   "documentType": <document type as a string>,
   "hawb": <HAWB number as a string>,
   "issuedOn": <Date of the document as RFC3339>,
   "pol": <Port of loading>,
   "poa": <Port of arrival>,
   "etd": <Date of departure as RFC3339>,
   "eta": <Date of arrival as RFC3339>,
   "flightNo": <Flight number as string>,
   "carrierName": <Unique name of a carrier company as a string>,
   "carrierAddress": <Address of a carrier company as a string>,
   "cargoAgentCode": <cargo agent code as a string as string>,"
   "shipperName": <Unique name of a shipper as a string>,
   "shipperAddress": <Address name of a shipper as a string>,
   "consigneeName": <Unique name of a consignee as a string>,
   "consigneeAddress": <Unique name of a consignee as a string>,
   "totalGrossWeight": <Total gross weight of the goods, in kg as a string>,
   "handlingInstructions": <handling instructions, information about how to transport the goods>
   "totalDimensions": {
       "length": <length of the load, in cm as a string>,
       "width": <width of the load, in cm as a string>,
       "height": <height of the load, in cm as a string>,
       "unit": <measurement unit as a string>
   },
   "shipmentOfPieces" : [
       {
            "itemNumber": <unique number or identifier of the item as string>,
            "itemDescription": <description of transported item as a string>,
            "quantity": <number of pieces as an integer>,
            "cartons": <number of cartons as string>,
            "weight": <total weight of items as string>,
            "unit": <measurement unit as string>,
            "hsCode": <customs clearance code as string>,
            "manufacturer": <unique name of a manufacturer>
       }
   ]
}
**Data Format Handling:**

* Possible aliases for "item": "Purchase Order", "PO", "Style", "Article"
* Length, width, and height can be presented in the format "<length> / <width> / <height>" or "<length> x <width> x <height>"
* Look for a number of cartons before "CTN" or "CTNS" or "CARTONS"
* Do not distinguish between upper and lower case for parsing the input
* Convert "pol" and "poa" to the IATA-Code
* Remove null values from the json to keep it compact
* If you decide that the document is not a HAWB, return only the fields "isHawb" and "documentType" with the value "false", or the document type you think fits respectively, as JSON.
* If the given language is not english, then translate corresponding fields to english prior to generating the JSON output.

**Error Handling:**

If any of the required fields ("hawb", "totalGrossWeight", "volume", and "cargoName") cannot be extracted from the PDF, generate an error message in the following format: 
	"ERROR: The following required fields could not be extracted from the PDF: [list of missing fields]. Please review the PDF and provide the missing information." 
For example, if "hawb" and "volume" are missing, the error message should be:
	"ERROR: The following required fields could not be extracted from the PDF: ["hawb", "volume"]. Please review the PDF and provide the missing information."
If all required fields are present, generate the JSON output as specified above.
For optional fields that are not found in the PDF, dont include them in the JSON output.
Prioritize accuracy and ensure the generated JSON is valid and conforms to the specified structure.`

type HwbService struct {
	config *configuration.Config
	ai     *AiCommunicationService
}

func NewHwbService(config *configuration.Config) *HwbService {
	ai := NewAiCommunicationService(config)
	return &HwbService{config: config, ai: ai}
}

func (i *HwbService) AnalysePdfFile(filename string) error {
	answer, err := i.ai.GenerateContentFromPDF(filename, promptHwbVertexAI)
	if err != nil {
		return err
	}

	cleanedString := strings.TrimPrefix(answer, "```json")
	answer = strings.TrimSuffix(cleanedString, "```")
	if os.Getenv("WRITE_TO_FILE") == "true" {
		err = writeToFile(answer, filename)
		if err == nil {
			fmt.Println("created file: " + filename)
		}
	}
	return err
}

func writeToFile(answer string, fileName string) error {
	create, _ := os.Create("./ai-output/" + fileName + ".json")
	_, _ = create.WriteString(answer)
	return create.Close()
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
	factoryName := responseVertexAI.ShipperName
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
		}
		if factoryName != "" {
			product.Manufacturer = iata.Manufacturer{
				Type: "cargo:Company",
				Name: factoryName,
			}
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
