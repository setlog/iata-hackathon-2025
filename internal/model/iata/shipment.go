package iata

type Shipment struct {
	Context          Context         `json:"@context"`
	Type             string          `json:"@type"`
	GoodsDescription string          `json:"cargo:goodsDescription"`
	TotalGrossWeight Measurement     `json:"cargo:totalGrossWeight"`
	TotalDimensions  TotalDimensions `json:"cargo:totalDimensions"`
	ShipmentOfPieces []ObjectLink    `json:"cargo:shipmentOfPieces"`
	InvolvedParties  []InvolvedParty `json:"cargo:involvedParties"`
}

type TotalDimensions struct {
	Height Measurement `json:"cargo:height"`
	Length Measurement `json:"cargo:length"`
	Width  Measurement `json:"cargo:width"`
	Volume Measurement `json:"cargo:volume"`
}
