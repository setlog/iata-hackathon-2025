package iata

type Person struct {
	Context     Context `json:"@context"`
	Type        string  `json:"@type"`
	Salutation  string  `json:"cargo:salutation"`
	LastName    string  `json:"cargo:lastName"`
	FirstName   string  `json:"cargo:firstName"`
	ContactRole string  `json:"cargo:contactRole"`
}
