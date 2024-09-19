package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
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

var logger *slog.Logger = nil

func main() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	r := mux.NewRouter()
	r.HandleFunc("/testvalidation", handlerFunc)

	port := 8080
	logger.Info("Service has started", slog.Int("port", port))

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func handlerFunc(w http.ResponseWriter, req *http.Request) {

	fileName := req.URL.Query().Get("fileName")
	if fileName == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("parameter 'fileName' not found.")
		fmt.Fprint(w, "parameter 'fileName' not found.")
		return
	}

	responseVertexAI, err := generateContentFromPDF(fileName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("%v", slog.Any("error", err))
		fmt.Fprintf(w, "%v", err)
		return
	}

	testValidationResp := convertResponse(*responseVertexAI)
	answer, err := json.Marshal(testValidationResp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("%v", slog.Any("error", err))
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
	bucketName := os.Getenv("GCLOUD_BUCKETNAME")

	if err := waitUntilFileExists(bucketName, fileName); err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)

	part := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  fmt.Sprintf("%s/%s", bucketName, fileName),
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

func waitUntilFileExists(bucketName, fileName string) error {
	bucketName, _ = strings.CutPrefix(bucketName, "gs://")
	parts := strings.SplitN(bucketName, "/", 2)
	if len(parts) > 1 {
		bucketName = parts[0]
		fileName = parts[1] + "/" + fileName
	}
	timeOutInSeconds := 60
	for i := 0; i < timeOutInSeconds; i++ {
		bExists, err := isFileExist(bucketName, fileName)
		if bExists {
			return nil
		}
		if err != nil {
			logger.Warn("cannot read the cloud file", slog.String("bucket", bucketName), slog.String("file", fileName), slog.Any("error", err))
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout occured while reading the bucket %s, file %s", bucketName, fileName)
}

func isFileExist(bucketName, fileName string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	// Get the object handle
	obj := client.Bucket(bucketName).Object(fileName)

	// Try to get the attributes of the object
	_, err = obj.Attrs(ctx)
	if err != nil {
		// If the error is a "NotFound" error, the file doesn't exist
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		// Other errors should be handled accordingly
		return false, fmt.Errorf("failed to get object attributes: %v", err)
	}

	// The object exists
	return true, nil
}
