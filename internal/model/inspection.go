package model

type InspectionResponseVertexAI struct {
	InspectionReportNo   string `json:"InspectionReportNo"`
	InspectionDate       string `json:"InspectionDate"`
	InspectionDateFormat string `json:"InspectionDateFormat"`
	InspectionResult     string `json:"InspectionResult"`
	Inspector            string `json:"Inspector"`
}

type InspectionTestValidationResponse struct {
	InspectionReportNo   string `json:"InspectionReportNo"`
	InspectionDate       string `json:"InspectionDate"`
	InspectionDateFormat string `json:"InspectionDateFormat"`
	InspectionResult     string `json:"InspectionResult"`
	Inspector            string `json:"Inspector"`
}
