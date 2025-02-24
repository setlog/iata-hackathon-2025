package iata

type Hwb struct {
	Context         Context         `json:"@context"`
	Type            string          `json:"@type"`
	WaybillNumber   string          `json:"cargo:waybillNumber"`
	WaybillType     string          `json:"cargo:waybillType"`
	Shipment        ObjectLink      `json:"cargo:shipment"`
	InvolvedParties []InvolvedParty `json:"cargo:involvedParties"`
}
