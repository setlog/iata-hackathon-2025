package main

type ResponseItemVertexAI struct {
	Name   string `json:"Test"`
	Status string `json:"Status"`
}

type ResponseVertexAI struct {
	Type            string                 `json:"DocumentType"`
	Status          string                 `json:"Status"`
	Summary         string                 `json:"Summary"`
	Results         []ResponseItemVertexAI `json:"Results"`
	ReasonOfFailure string                 `json:"ReasonOfFailure"`
	TestReportNo    string                 `json:"TestReportNo"`
	TesteDate       string                 `json:"TestDate"`
	DateFormatUsed  string                 `json:"DateFormatUsed"`
	Lab             string                 `json:"Lab"`
}

type TestValidationItem struct {
	Test   string `json:"Test"`
	Status string `json:"ItemStatus"`
}

type TestValidationResponse struct {
	DocumentType    string               `json:"DocumentType"`
	Status          string               `json:"Status"`
	ReasonOfFailure string               `json:"ReasonOfFailure"`
	Items           []TestValidationItem `json:"Items"`
	TestReportNo    string               `json:"TestReportNo"`
	TesteDate       string               `json:"TestDate"`
	DateFormatUsed  string               `json:"DateFormatUsed"`
	Lab             string               `json:"Lab"`
}
