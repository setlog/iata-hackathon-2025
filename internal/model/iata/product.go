package iata

type product struct {
	Context          Context      `json:"@context"`
	Type             string       `json:"@type"`
	Description      string       `json:"cargo:description"`
	HsCode           string       `json:"cargo:hsCode"`
	UniqueIdentifier string       `json:"cargo:uniqueIdentifier"`
	Manufacturer     Manufacturer `json:"cargo:manufacturer"`
}

type Manufacturer struct {
	Type string `json:"@type"`
	Name string `json:"cargo:name"`
}
