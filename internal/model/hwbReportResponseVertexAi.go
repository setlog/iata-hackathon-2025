package model

type HwbReportResponseVertexAi struct {
	Hawb                 string          `json:"hawb"`
	Issued               string          `json:"issued"`
	Pol                  string          `json:"pol"`
	Poa                  string          `json:"poa"`
	Etd                  string          `json:"etd"`
	Eta                  string          `json:"eta"`
	Flightno             string          `json:"flightno"`
	CarrierName          string          `json:"carrierName"`
	CarrierAdress        string          `json:"carrierAdress"`
	ShipperName          string          `json:"shipperName"`
	ShippperAdress       string          `json:"shippperAdress"`
	ConsigneeName        string          `json:"consigneeName"`
	ConsigneeAdress      string          `json:"consigneeAdress"`
	TotalGrossWeight     string          `json:"totalGrossWeight"`
	HandlingInstructions string          `json:"handlingInstructions"`
	TotalDimensions      TotalDimensions `json:"totalDimensions"`
	ShipmentOfPieces     []Pieces        `json:"shipmentOfPieces"`
	FactoryName          string          `json:"factoryName"`
	FactoryAdress        string          `json:"factoryAdress"`
}

type TotalDimensions struct {
	Length string `json:"length"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Volume string `json:"volume"`
}

type Pieces struct {
	ItemNumber      string `json:"itemNumber"`
	ItemDescription string `json:"itemDescription"`
	Quantity        int    `json:"quantity"`
	Cartons         string `json:"cartons"`
	Weight          string `json:"weight"`
	Unit            string `json:"unit"`
	HsCode          string `json:"hsCode"`
	Manufacturer    string `json:"manufacturer"`
}
