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
            Please summarize the given document and report me, if the test results in this document show a passed or failed status.
            Give me a feedback as plain json (no "json" encasing) with fields "Status", "Summary", a list of tested product features in "Results" and the reason of failure in "ReasonOfFailure", if it failed
Here is the remplate for the response json struct
{
    "Status": "PASS" | "FAIL",
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
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handlerFunc(w http.ResponseWriter, req *http.Request) {

	fileName := req.URL.Query().Get("fileName")
	if fileName == "" {
		fmt.Fprint(w, "parameter 'fileName' not found.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseVertexAI, err := generateContentFromPDF(fileName)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	testValidationResp := convertResponse(*responseVertexAI)
	answer, err := json.Marshal(testValidationResp)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", answer)
	w.WriteHeader(http.StatusOK)
}

func convertResponse(responseVertexAI ResponseVertexAI) TestValidationResponse {
	return TestValidationResponse{
		Status:          responseVertexAI.Status,
		ReasonOfFailure: responseVertexAI.ReasonOfFailure,
	}
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
