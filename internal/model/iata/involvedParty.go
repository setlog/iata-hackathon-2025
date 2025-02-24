package iata

type InvolvedParty struct {
	Type         string     `json:"@type"`
	Role         string     `json:"cargo:role"`
	Organization ObjectLink `json:"cargo:organization"`
}
