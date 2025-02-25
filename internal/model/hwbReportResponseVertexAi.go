package model

type HwbReportResponseVertexAi struct {
	IsHawb               bool            `json:"isHawb"`
	DocumentType         string          `json:"documentType,omitempty"`
	Hawb                 string          `json:"hawb,omitempty"`
	Issued               string          `json:"issued,omitempty"`
	Pol                  string          `json:"pol,omitempty"`
	Poa                  string          `json:"poa,omitempty"`
	Etd                  string          `json:"etd,omitempty"`
	Eta                  string          `json:"eta,omitempty"`
	Flightno             string          `json:"flightno,omitempty"`
	CarrierName          string          `json:"carrierName,omitempty"`
	CarrierAdress        string          `json:"carrierAdress,omitempty"`
	ShipperName          string          `json:"shipperName,omitempty"`
	ShippperAdress       string          `json:"shippperAdress,omitempty"`
	ConsigneeName        string          `json:"consigneeName,omitempty"`
	ConsigneeAdress      string          `json:"consigneeAdress,omitempty"`
	TotalGrossWeight     string          `json:"totalGrossWeight,omitempty"`
	HandlingInstructions string          `json:"handlingInstructions,omitempty"`
	TotalDimensions      TotalDimensions `json:"totalDimensions,omitempty"`
	ShipmentOfPieces     []Pieces        `json:"shipmentOfPieces,omitempty"`
	FactoryName          string          `json:"factoryName,omitempty"`
	FactoryAdress        string          `json:"factoryAdress,omitempty"`
}

type TotalDimensions struct {
	Length string `json:"length,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

type Pieces struct {
	ItemNumber      string `json:"itemNumber,omitempty"`
	ItemDescription string `json:"itemDescription,omitempty"`
	Quantity        int    `json:"quantity,omitempty"`
	Cartons         string `json:"cartons,omitempty"`
	Weight          string `json:"weight,omitempty"`
	Unit            string `json:"unit,omitempty"`
	HsCode          string `json:"hsCode,omitempty"`
	Manufacturer    string `json:"manufacturer,omitempty"`
}
