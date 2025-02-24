package iata

type Item struct {
	Context            Context     `json:"@context"`
	Type               string      `json:"@type"`
	Measurement        Measurement `json:"cargo:itemQuantity"`
	DescribedByProduct ObjectLink  `json:"cargo:describedByProduct"`
}
