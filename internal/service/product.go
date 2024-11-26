package service

import (
	"encoding/json"

	"com.setlog/internal/configuration"
	"com.setlog/internal/model"
)

const promptProductTestVertexAI = `You are a very professional specialist in quality control of the consumer goods.
            Please check if the given document is a report of a consumer good quality control test.
            If it is one, then please summarize the given document and report me, if the test results in this document show a passed or failed status.
            If it is not, then please set the Status to "UNDEFINED", "DocumentType" to the corresponding description and the rest of the fields to null.
			Also provide me general document information, i want know the test report number "TestReportNo", the test date "TestDate" and the lab/institute "Lab" which performed the product test.
            Give me a feedback as plain json (no json encasing) with fields "Status", "Summary", a list of tested product features in "Results" and the reason of failure in "ReasonOfFailure", if it failed
Here is the remplate for the response json struct
{
    "DocumentType": "REPORT" | "<string with the description of what the document is>",
    "Status": "PASS" | "FAIL" | "UNDEFINED",
    "Summary": "<summary text>",
    "Items": [
        {
            "Test": "<name of the tested product feature>",
            "ItemStatus": "PASS" | "FAIL"
        },
        ...
    ],
    "ReasonOfFailure": "<description of the reason the test failed, empty string if the test passed>"
	"TestReportNo": "222416610",
	"TestDate": "27-05-2024",
	"DateFormatUsed": "DD.MM.YYYY",
	"Lab": "TÜV Rheinland"
}`

type ProductTestService struct {
	config *configuration.Config
	ai     *AiCommunicationService
}

func NewProductTestService(config *configuration.Config) *ProductTestService {
	ai := NewAiCommunicationService(config)
	return &ProductTestService{config: config, ai: ai}
}

func (p *ProductTestService) AnalysePdfFile(filename string) ([]byte, error) {

	answer, err := p.ai.GenerateContentFromPDF(filename, promptProductTestVertexAI)
	result := model.ProductTestResponseVertexAI{}
	err = json.Unmarshal([]byte(answer), &result)
	if err != nil {
		return nil, err
	}
	ValidationResp := p.convertResponse(&result)
	return json.Marshal(ValidationResp)
}

func (p *ProductTestService) convertResponse(responseVertexAI *model.ProductTestResponseVertexAI) model.ProductTestValidationResponse {
	result := model.ProductTestValidationResponse{
		DocumentType:    responseVertexAI.Type,
		Status:          responseVertexAI.Status,
		ReasonOfFailure: responseVertexAI.ReasonOfFailure,
		TestReportNo:    responseVertexAI.TestReportNo,
		TesteDate:       responseVertexAI.TesteDate,
		DateFormatUsed:  responseVertexAI.DateFormatUsed,
		Lab:             responseVertexAI.Lab,
	}
	for _, t := range responseVertexAI.Results {
		result.Items = append(result.Items, model.ProductTestValidationItem{
			Test:   t.Name,
			Status: t.Status,
		})
	}
	return result
}
