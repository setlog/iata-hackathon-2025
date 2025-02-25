# HAWB to IATA One Record Converter

This Golang project parses unstructured HAWB documents in PDF format and converts extracted data into the [IATA One Record API](https://www.iata.org/en/programs/cargo/one-record/) using [Vertex AI](https://cloud.google.com/vertex-ai) as an LLM.

## Prerequisites

1. **Golang**: Ensure that [Go](https://go.dev/doc/install) is installed on your system.
2. **Google Cloud Account**: The project uses Google Cloud services such as Vertex AI and Cloud Storage.
3. **.env File**: Create a `.env` file with the required configuration values.

## Configuration

Create a `.env` file in the root directory of the project with the following content:

```ini
GOOGLE_APPLICATION_CREDENTIALS=<Path to credentials on Google Cloud>
GCLOUD_PROJECT_ID=<ID of Google Cloud Project>
GCLOUD_LOCATION=<Location of Google Cloud Project>
GCLOUD_BUCKETNAME=<Path to bucket to use on Google Cloud Project>
AI_MODEL=<AI model for document parsing, e.g. gemini-2.0-flash-exp>
```

## Installation
1. Clone this repository:
```bash
git clone <repository-url>
cd haweb2iata-converter
```
2. Install dependencies:
```bash
go mod tidy
```

## Running the Application
Start the service using the following command:
```bash
go run -o haweb2iata-converter .
```

## API Endpoints
The application provides the following REST API endpoints:
### 1. PDF Analysis
- `GET /hwbreportanalysis/all`
  - Description: Parses all PDF documents (HAWBs) stored in the Google Cloud Storage bucket GCLOUD_BUCKETNAME and returns the output in JSON format.
- `GET /hwbreportanalysis?fileName=<fileName>`
  - Description: Parses the specified PDF document <fileName> from the bucket GCLOUD_BUCKETNAME and returns the output in JSON format.
### 2. Convert JSON to IATA One Record API
- `GET /json2iata/all`
  - Description: Converts all generated JSON data into IATA format and uploads it to the One Record Server via API.
- `POST /json2iata`
  - Description: Converts a single HAWB JSON into IATA format and uploads it to the One Record Server via API.
  - Request Body: JSON containing the data to be converted.

## Contact
For any questions or issues, please contact the project maintainer or create an Issue.
