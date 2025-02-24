package service

import (
	"encoding/json"

	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
)

const promptHwbVertexAI = `You are a very professional specialist in analysing inspections of consumer goods.
            Please check if the given document is a report of a consumer good inspection test.
            If it is one, then please summarize the given document and report me, if the test results in this document show a passed or failed status.
            If it is not, then please set the Status to "UNDEFINED" and the rest of the fields to null.
			Also provide me general document information, i want know the inspection report number "InspectionReportNo", the inspection date "InspectionDate" and the format of the date e.g. yyyy-mm-dd "InspectionDateFormat".
			Also give me the person who performed the inspection. This field is most often marked with "Inspected by" in the inspection report. The name of the json field is "Inspector".
            Give me a feedback as plain json (no json encasing) with fields "InspectionResult", "InspectionDate" and "InspectionReportNo"
Here is the remplate for the response json struct
{
    "InspectionResult": "PASS" | "FAIL" | "UNDEFINED",
	"InspectionReportNo": "222416610",
	"InspectionDate": "27-05-2024",
	"InspectionDateFormat": "DD.MM.YYYY",
	"Inspector": "Spark Ling"
}`

type HwbService struct {
	config *configuration.Config
	ai     *AiCommunicationService
}

func NewHwbService(config *configuration.Config) *HwbService {
	ai := NewAiCommunicationService(config)
	return &HwbService{config: config, ai: ai}
}

func (i *HwbService) AnalysePdfFile(filename string) (*model.InspectionTestValidationResponse, error) {
	answer, err := i.ai.GenerateContentFromPDF(filename, promptHwbVertexAI)
	result := model.HwbReportResponseVertexAi{}
	err = json.Unmarshal([]byte(answer), &result)
	if err != nil {
		return nil, err
	}
	ValidationResp := i.convertResponse(&result)
	return ValidationResp, nil
}

func (i *HwbService) convertResponse(responseVertexAI *model.HwbReportResponseVertexAi) *model.InspectionTestValidationResponse {

	return nil
}
