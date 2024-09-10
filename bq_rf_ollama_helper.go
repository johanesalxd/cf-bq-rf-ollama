package bqrfollama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// textsToTexts processes multiple text inputs concurrently using via Ollama API
func textsToTexts(ctx context.Context, bqReq *BigQueryRequest) *BigQueryResponse {
	// Initialize a slice to store the processed texts
	texts := make([]string, len(bqReq.Calls))
	wg := new(sync.WaitGroup)
	semaphore := make(chan struct{}, concurrencyLimit)

	for i, call := range bqReq.Calls {
		select {
		case <-ctx.Done():
			log.Printf("Context cancelled before starting goroutine #%d", i)
			texts[i] = string(GenerateJSONResponse(json.RawMessage(`{"error":"Request cancelled"}`)))

			continue
		default:
			wg.Add(1)

			// Process each call concurrently
			go func(i int, call []interface{}) {
				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() {
					// Release semaphore
					<-semaphore
					wg.Done()
				}()
				log.Printf("Processing request in Goroutine #%d", i)

				// Check if call has 2 elements
				if len(call) != 2 {
					log.Printf("Error in goroutine #%d: call does not have enough elements", i)
					texts[i] = string(GenerateJSONResponse(json.RawMessage(`{"error":"Invalid input: expected 2 elements"}`)))

					return
				}

				// Update the input from the call slice
				input := OllamaRequest{
					Prompt: fmt.Sprint(call[0]),
					Model:  fmt.Sprint(call[1]),
					// TODO: Implement streaming
					Stream: false,
				}

				// Send the request to Ollama and handle the response
				output, err := SendOllamaRequest(ctx, input)
				if err != nil {
					log.Printf("Error in goroutine #%d generating text for input: %v", i, err)
					texts[i] = string(GenerateJSONResponse(json.RawMessage(fmt.Sprintf(`{"error":"%s"}`, err.Error()))))

					return
				}

				texts[i] = string(GenerateJSONResponse(output))
			}(i, call)
		}
	}
	wg.Wait()

	// Prepare and return the BigQuery response
	return &BigQueryResponse{
		Replies: texts,
	}
}

// GenerateJSONResponse converts the input to JSON format
func GenerateJSONResponse(input any) json.RawMessage {
	jsonInput, err := json.Marshal(input)
	if err != nil {
		log.Printf("Error marshaling input to JSON: %v", err)
		return json.RawMessage(fmt.Sprintf(`{"error":"%s"}`, err.Error()))
	}

	return jsonInput
}
