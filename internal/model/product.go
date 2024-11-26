package model

type ProductTestResponseItemVertexAI struct {
	Name   string `json:"Test"`
	Status string `json:"ItemStatus"`
}

type ProductTestResponseVertexAI struct {
	Type            string                            `json:"DocumentType"`
	Status          string                            `json:"Status"`
	Summary         string                            `json:"Summary"`
	Items           []ProductTestResponseItemVertexAI `json:"Items"`
	ReasonOfFailure string                            `json:"ReasonOfFailure"`
	TestReportNo    string                            `json:"TestReportNo"`
	TesteDate       string                            `json:"TestDate"`
	DateFormatUsed  string                            `json:"DateFormatUsed"`
	Lab             string                            `json:"Lab"`
}

type ProductTestValidationItem struct {
	Test       string `json:"Test"`
	ItemStatus string `json:"ItemStatus"`
}

type ProductTestValidationResponse struct {
	DocumentType    string                      `json:"DocumentType"`
	Status          string                      `json:"Status"`
	ReasonOfFailure string                      `json:"ReasonOfFailure"`
	Items           []ProductTestValidationItem `json:"Items"`
	TestReportNo    string                      `json:"TestReportNo"`
	TesteDate       string                      `json:"TestDate"`
	DateFormatUsed  string                      `json:"DateFormatUsed"`
	Lab             string                      `json:"Lab"`
}
