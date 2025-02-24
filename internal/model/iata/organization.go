package iata

type Organization struct {
	Context        Context    `json:"@context"`
	Type           string     `json:"@type"`
	Name           string     `json:"cargo:name"`
	ContactPersons ObjectLink `json:"cargo:contactPersons"`
}
