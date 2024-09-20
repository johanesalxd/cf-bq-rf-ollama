package bqrfollama

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// BQRFOllama handles HTTP requests for the BigQuery Remote Function using Ollama AI
func BQRFOllama(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming BigQuery request
	bqReq := new(BigQueryRequest)
	if err := json.NewDecoder(r.Body).Decode(bqReq); err != nil {
		SendError(w, err, http.StatusBadRequest)

		return
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(contextTimeoutS)*time.Second)
	defer func() {
		cancel()
		log.Printf("Done, Goroutines closed due to: %v", ctx.Err())
	}()

	// Process the request using ollama
	bqResp := textsToTexts(ctx, bqReq)

	// Send the successful response
	SendSuccess(w, bqResp)
}
