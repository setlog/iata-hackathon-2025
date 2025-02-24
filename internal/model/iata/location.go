package iata

type Location struct {
	Type         string `json:"@type"`
	LocationType string `json:"cargo:locationType"`
	Code         string `json:"cargo:code"`
}
