package iata

type TransportMeans struct {
	Context               Context    `json:"@context"`
	Type                  string     `json:"@type"`
	VehicleModel          string     `json:"cargo:vehicleModel"`
	VehicleRegistration   string     `json:"cargo:vehicleRegistration"`
	TransportOrganization ObjectLink `json:"cargo:transportOrganization"`
}
