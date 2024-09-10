package bqrfollama_test

import (
	"encoding/json"
	"testing"

	bqrfollama "github.com/johanesalxd/cf-bq-rf-ollama"
)

// TestGenerateJSONResponse tests the GenerateJSONResponse function
func TestGenerateJSONResponse(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// Test case 1: Ensure a valid struct is correctly serialized
		{
			name:     "Valid struct",
			input:    struct{ Name string }{"John"},
			expected: `{"Name":"John"}`,
		},
		// Test case 2: Ensure a valid map is correctly serialized
		{
			name:     "Valid map",
			input:    map[string]int{"Age": 30},
			expected: `{"Age":30}`,
		},
		// Test case 3: Ensure nil input is handled correctly
		{
			name:     "Nil input",
			input:    nil,
			expected: `null`,
		},
		// Test case 4: Ensure unmarshalable input returns an error message
		{
			name:     "Unmarshalable input",
			input:    make(chan int),
			expected: `{"error":"json: unsupported type: chan int"}`,
		},
		// Test case 5: Ensure error input is handled correctly
		{
			name:     "Error input",
			input:    json.RawMessage(`{"error":"Request cancelled"}`),
			expected: `{"error":"Request cancelled"}`,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bqrfollama.GenerateJSONResponse(tt.input)

			// Compare JSON strings
			if string(result) != tt.expected {
				t.Errorf("generateJSONResponse() = %v, want %v", string(result), tt.expected)
			}

			// Verify it's valid JSON
			var js json.RawMessage
			if err := json.Unmarshal(result, &js); err != nil {
				t.Errorf("generateJSONResponse() produced invalid JSON: %v", err)
			}
		})
	}
}
