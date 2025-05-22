# Go Lambda API Fetcher
A production-ready AWS Lambda function written in Go that fetches data from a public JSON API, logs the total number of items, and prints the title of the first item—all with proper error handling and logging.

## Overview
This serverless function demonstrates best practices for creating an AWS Lambda in Go, including:

- Robust error handling with retries and exponential backoff

- Comprehensive logging

- Configurable API endpoints via environment variables

- AWS Lambda integration

- Proper response formatting

- Unit testing implementation

## Requirements
- Go 1.18 or later

- AWS CLI (for deployment)

- AWS account with appropriate permissions

Project Structure

lambda-api-fetcher/

├── main.go              # Main Lambda handler code

├── main_test.go         # Unit tests

├── go.mod               # Go module definition

├── go.sum               # Go module checksums

└── README.md            # Project documentation

## Installation

1. Clone the repository:
```
git clone https://github.com/sweetdream-001/lambda-api-fetcher.git
cd lambda-api-fetcher
```
2. Install dependencies:
```sh
go mod download
```
## Local Development
### Running Locally
You can run the code locally to test its functionality:

```bash
go run main.go
```
Note: When running locally, the code will execute but immediately exit as it's designed as a Lambda function.

### Building
To build the Lambda deployment package:

```bash
GOOS=linux GOARCH=amd64 go build -o main
zip function.zip main
```
## Configuration
The Lambda function supports the following environment variables:
```
Variable	           Description	                                   Default
API_URL	         URL of the JSON API to fetch	     https://jsonplaceholder.typicode.com/posts
MAX_RETRIES	     Maximum number of retry attempts	   3
TIMEOUT_SECONDS	 HTTP client timeout in seconds	       10
```
## Deployment
### First-time Deployment
```bash
aws lambda create-function \
  --function-name GoAPIFetcher \
  --runtime go1.x \
  --handler main \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::YOUR_ACCOUNT_ID:role/lambda-execution-role \
  --environment Variables={API_URL=https://jsonplaceholder.typicode.com/posts,MAX_RETRIES=3}
  ```
### Updating the Function
```bash
aws lambda update-function-code \
  --function-name GoAPIFetcher \
  --zip-file fileb://function.zip
  ```
### Updating Configuration
``` bash
aws lambda update-function-configuration \
  --function-name GoAPIFetcher \
  --environment Variables={API_URL=https://api.example.com/data,MAX_RETRIES=5}
  ```
## Testing
### Running Unit Tests
```bash
go test -v
```
### Testing in AWS Console
- Navigate to the Lambda function in the AWS Console

- Click on the "Test" tab

- Create a new test event with an empty JSON object {}

- Click "Test" to execute the function

## Error Handling
The function implements several layers of error handling:

- HTTP client timeout configuration

- Context-based cancellation

- Response status code validation

- Retry mechanism with exponential backoff

- Comprehensive error logging

- Graceful failure responses

## Lambda Response Format
### Successful response:

```json
{
  "statusCode": 200,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "{\"total\":100,\"firstItemTitle\":\"Example title\"}"
}
```
### Error response:

```json
{
  "statusCode": 500,
  "headers": {
    "Content-Type": "application/json"
  },
  "body": "{\"error\":\"Failed to fetch posts: error details\"}"
}
```
### Security Considerations
- The function uses HTTPS for API communication

- No sensitive data is logged

- Environment variables are used for configuration

- Proper input validation is implemented

- Timeouts prevent hanging connections

## Contributing
- Fork the repository

- Create a feature branch (git checkout -b feature/amazing-feature)

- Commit your changes (git commit -m 'Add some amazing feature')

- Push to the branch (git push origin feature/amazing-feature)

- Open a Pull Request

### I'm excited about the opportunity at SERASAR. 
### With prior experience at the company and a strong developer skillset, 
### I'm confident I can contribute effectively to your team.


Acknowledgements
- AWS Lambda Go Runtime

- JSONPlaceholder for providing a free test API