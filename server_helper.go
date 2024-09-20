package bqrfollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	initOnce                          sync.Once
	concurrencyLimit, contextTimeoutS int
	httpClient                        *http.Client
)

// SendError sends an error response with the given error message and HTTP status code
func SendError(w http.ResponseWriter, err error, code int) {
	bqResp := new(BigQueryResponse)
	bqResp.ErrorMessage = fmt.Sprintf("Got error with details: %v", err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(bqResp)
}

// SendSuccess sends a successful response with the given BigQueryResponse
func SendSuccess(w http.ResponseWriter, bqResp *BigQueryResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bqResp)
}

// TODO: change from HTTP call to library (https://github.com/ollama/ollama/blob/main/examples/go-generate)
func sendOllamaRequest(ctx context.Context, req OllamaRequest) (json.RawMessage, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	url := fmt.Sprintf("%s/api/generate", os.Getenv("OLLAMA_URL"))
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return body, nil
}

func initAll() {
	var err error

	// Parse the concurrency limit from the environment variable
	concurrencyLimit, err = strconv.Atoi(os.Getenv("CONCURRENCY_LIMIT"))
	if err != nil {
		log.Printf("Failed to parse CONCURRENCY_LIMIT, using default value of 100: %v", err)
		concurrencyLimit = 100
	}

	// Parse the context timeout from the environment variable
	contextTimeoutS, err = strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT_S"))
	if err != nil {
		log.Printf("Failed to parse CONTEXT_TIMEOUT_S, using default value of 30 seconds: %v", err)
		contextTimeoutS = 30
	}

	// Initialize the HTTP client
	httpClient = &http.Client{}
}
