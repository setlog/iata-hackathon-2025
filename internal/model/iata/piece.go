package iata

type Piece struct {
	Context              Context               `json:"@context"`
	Type                 string                `json:"@type"`
	Coload               bool                  `json:"cargo:coload"`
	GoodsDescription     string                `json:"cargo:goodsDescription"`
	Upid                 string                `json:"cargo:upid"`
	ContainedItems       []ObjectLink          `json:"cargo:containedItems"`
	HandlingInstructions []HandlingInstruction `json:"cargo:handlingInstructions"`
}

type HandlingInstruction struct {
	Type               string `json:"@type"`
	ServiceType        string `json:"cargo:serviceType"`
	ServiceDescription string `json:"cargo:serviceDescription"`
	ServiceTypeCode    string `json:"cargo:serviceTypeCode"`
}
