package iata

type Measurement struct {
	Type           string  `json:"@type"`
	Unit           string  `json:"cargo:unit"`
	NumericalValue float64 `json:"cargo:numericalValue"`
}
