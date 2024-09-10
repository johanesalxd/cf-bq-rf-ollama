package bqrfollama_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	bqrfollama "github.com/johanesalxd/cf-bq-rf-ollama"
	"github.com/stretchr/testify/assert"
)

func TestSendError(t *testing.T) {
	// Test cases for different error scenarios
	tests := []struct {
		name string
		err  error
		code int
		want string
	}{
		{
			name: "error with code 500",
			err:  fmt.Errorf("internal server error"),
			code: http.StatusInternalServerError,
			want: "Got error with details: internal server error",
		},
		{
			name: "error with code 400",
			err:  fmt.Errorf("bad request"),
			code: http.StatusBadRequest,
			want: "Got error with details: bad request",
		},
		{
			name: "error with code 404",
			err:  fmt.Errorf("not found"),
			code: http.StatusNotFound,
			want: "Got error with details: not found",
		},
		{
			name: "error with code 403",
			err:  fmt.Errorf("forbidden"),
			code: http.StatusForbidden,
			want: "Got error with details: forbidden",
		},
		{
			name: "error with code 401",
			err:  fmt.Errorf("unauthorized"),
			code: http.StatusUnauthorized,
			want: "Got error with details: unauthorized",
		},
		{
			name: "error with custom message",
			err:  fmt.Errorf("custom error message"),
			code: http.StatusInternalServerError,
			want: "Got error with details: custom error message",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP response recorder
			w := httptest.NewRecorder()
			// Call the SendError function with test case parameters
			bqrfollama.SendError(w, test.err, test.code)

			// Get the response from the recorder
			resp := w.Result()
			// Check if the status code matches the expected code
			if resp.StatusCode != test.code {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, resp.StatusCode, test.code)
			}
			// Verify that the Content-Type header is set to application/json
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, resp.Header.Get("Content-Type"), "application/json")
			}

			// Decode the response body into a BigQueryResponse struct
			var got bqrfollama.BigQueryResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, err, test.want)
			}
			// Check if the error message in the response matches the expected message
			if got.ErrorMessage != test.want {
				t.Errorf("SendError(%v, %v) = %v, want %v", test.err, test.code, got.ErrorMessage, test.want)
			}
		})
	}
}

func TestSendSuccess(t *testing.T) {
	// Test cases for different success scenarios
	tests := []struct {
		name string
		resp *bqrfollama.BigQueryResponse
		want []string
	}{
		{
			name: "single success reply",
			resp: &bqrfollama.BigQueryResponse{
				Replies: []string{"success"},
			},
			want: []string{"success"},
		},
		{
			name: "multiple success replies",
			resp: &bqrfollama.BigQueryResponse{
				Replies: []string{"success1", "success2", "success3"},
			},
			want: []string{"success1", "success2", "success3"},
		},
		{
			name: "empty replies",
			resp: &bqrfollama.BigQueryResponse{
				Replies: []string{},
			},
			want: []string{},
		},
		{
			name: "nil replies",
			resp: &bqrfollama.BigQueryResponse{
				Replies: nil,
			},
			want: nil,
		},
		{
			name: "response with error message",
			resp: &bqrfollama.BigQueryResponse{
				Replies:      []string{"success"},
				ErrorMessage: "Some error occurred",
			},
			want: []string{"success"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new HTTP response recorder
			w := httptest.NewRecorder()
			// Call the SendSuccess function with the test case response
			bqrfollama.SendSuccess(w, test.resp)

			// Get the response from the recorder
			resp := w.Result()
			// Check if the status code is 200 OK
			if resp.StatusCode != http.StatusOK {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, resp.StatusCode, http.StatusOK)
			}
			// Verify that the Content-Type header is set to application/json
			if resp.Header.Get("Content-Type") != "application/json" {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, resp.Header.Get("Content-Type"), "application/json")
			}

			// Decode the response body into a BigQueryResponse struct
			var got bqrfollama.BigQueryResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, err, test.want)
			}
			// Check if the replies in the response match the expected replies
			if !reflect.DeepEqual(got.Replies, test.want) {
				t.Errorf("SendSuccess(%v) = %v, want %v", test.resp, got.Replies, test.want)
			}
			// Verify that the error message in the response matches the input error message
			if got.ErrorMessage != test.resp.ErrorMessage {
				t.Errorf("SendSuccess(%v) ErrorMessage = %v, want %v", test.resp, got.ErrorMessage, test.resp.ErrorMessage)
			}
		})
	}
}

func TestSendOllamaRequest(t *testing.T) {
	tests := []struct {
		name           string
		model          string
		prompt         string
		serverResponse string
		expectedError  bool
	}{
		{
			name:           "Successful request",
			model:          "test-model",
			prompt:         "test prompt",
			serverResponse: "mocked response",
			expectedError:  false,
		},
		{
			name:           "Server error",
			model:          "test-model",
			prompt:         "test prompt",
			serverResponse: "Internal Server Error",
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/generate", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var req bqrfollama.OllamaRequest
				json.NewDecoder(r.Body).Decode(&req)
				assert.Equal(t, tt.model, req.Model)
				assert.Equal(t, tt.prompt, req.Prompt)

				if tt.expectedError && tt.name == "Server error" {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(tt.serverResponse))
				} else {
					resp := bqrfollama.OllamaResponse{Response: tt.serverResponse}
					json.NewEncoder(w).Encode(resp)
				}
			}))
			defer server.Close()

			// Set environment variable
			os.Setenv("OLLAMA_URL", server.URL)
			defer os.Unsetenv("OLLAMA_URL")

			// Test
			req := bqrfollama.OllamaRequest{
				Model:  tt.model,
				Prompt: tt.prompt,
			}

			resp, err := bqrfollama.SendOllamaRequest(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.serverResponse, resp.Response)
			}
		})
	}
}
