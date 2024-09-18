package main

type ResponseItemVertexAI struct {
	Name   string `json:"Test"`
	Status string `json:"Status"`
}

type ResponseVertexAI struct {
	Status          string                 `json:"Status"`
	Summary         string                 `json:"Summary"`
	Results         []ResponseItemVertexAI `json:"Results"`
	ReasonOfFailure string                 `json:"ReasonOfFailure"`
}

type TestValidationResponse struct {
	Status          string `json:"Status"`
	ReasonOfFailure string `json:"ReasonOfFailure"`
}
