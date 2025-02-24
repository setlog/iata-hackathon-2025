package iata

type TransportMovement struct {
	Context             Context        `json:"@context"`
	Type                string         `json:"@type"`
	ArrivalLocation     Location       `json:"cargo:arrivalLocation"`
	DepartureLocation   Location       `json:"cargo:departureLocation"`
	TransportMeans      ObjectLink     `json:"cargo:transportMeans"`
	TransportIdentifier string         `json:"cargo:transportIdentifier"`
	CargoMovementTimes  []MovementTime `json:"cargo:cargoMovementTimes"`
}

type MovementTime struct {
	Type              string            `json:"@type"`
	MovementMilestone string            `json:"cargo:movementMilestone"`
	MovementTimestamp MovementTimestamp `json:"cargo:movementTimestamp"`
}

type MovementTimestamp struct {
	Type  string `json:"@type"`
	Value string `json:"@value"`
}
