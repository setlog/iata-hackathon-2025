package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/vertexai/genai"
	"github.com/gorilla/mux"
)

const promptVertexAI = `You are a very professional specialist in quality control of the consumer goods.
            Please check if the given document is a report of a consumer good quality control test.
            If it is one, then please summarize the given document and report me, if the test results in this document show a passed or failed status.
            If it is not, then please set the Status to "UNDEFINED", "DocumentType" to the corresponding description and the rest of the fields to null.
            Give me a feedback as plain json (no json encasing) with fields "Status", "Summary", a list of tested product features in "Results" and the reason of failure in "ReasonOfFailure", if it failed
Here is the remplate for the response json struct
{
    "DocumentType": "REPORT" | "<string with the description of what the document is>",
    "Status": "PASS" | "FAIL" | "UNDEFINED",
    "Summary": "<summary text>",
    "Results": [
        {
            "Test": "<name of the tested product feature>",
            "Status": "PASS" | "FAIL"
        },
        ...
    ],
    "ReasonOfFailure": "<description of the reason the test failed, empty string if the test passed>"
}`

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/testvalidation", handlerFunc)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handlerFunc(w http.ResponseWriter, req *http.Request) {

	fileName := req.URL.Query().Get("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "parameter 'fileName' not found.")
		return
	}

	responseVertexAI, err := generateContentFromPDF(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	testValidationResp := convertResponse(*responseVertexAI)
	answer, err := json.Marshal(testValidationResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(answer))
}

func convertResponse(responseVertexAI ResponseVertexAI) TestValidationResponse {
	result := TestValidationResponse{
		DocumentType:    responseVertexAI.Type,
		Status:          responseVertexAI.Status,
		ReasonOfFailure: responseVertexAI.ReasonOfFailure,
	}
	for _, t := range responseVertexAI.Results {
		result.Items = append(result.Items, TestValidationItem{
			Test:   t.Name,
			Status: t.Status,
		})
	}
	return result
}

// generateContentFromPDF generates a response into the provided io.Writer, based upon the PDF
func generateContentFromPDF(fileName string) (*ResponseVertexAI, error) {
	modelName := "gemini-1.5-flash-001"
	projectID := os.Getenv("GCLOUD_PROJECT_ID")
	location := os.Getenv("GCLOUD_LOCATION")

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  fmt.Sprintf("%s/%s", os.Getenv("GCLOUD_BUCKETNAME"), fileName),
	}

	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockOnlyHigh,
		},
	}

	res, err := model.GenerateContent(ctx, part, genai.Text(promptVertexAI))
	if err != nil {
		return nil, fmt.Errorf("unable to generate contents: %w", err)
	}

	if len(res.Candidates) == 0 ||
		len(res.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("empty response from model")
	}

	answer := fmt.Sprintf("%v", res.Candidates[0].Content.Parts[0])
	result := ResponseVertexAI{}
	err = json.Unmarshal([]byte(answer), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
