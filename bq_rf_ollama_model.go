package bqrfollama

import (
	"encoding/json"
)

type PromptRequest struct {
	PromptInput string `json:"promptInput"`
	Model       string `json:"model"`
	//TODO: implement ModelConfig
	PromptOutput json.RawMessage `json:"promptOutput"`
}

type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Suffix  string                 `json:"suffix,omitempty"`
	Stream  bool                   `json:"stream"`
	Format  string                 `json:"format,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

func newPromptRequest() PromptRequest {
	return PromptRequest{
		PromptInput:  "",
		Model:        "",
		PromptOutput: json.RawMessage(""),
	}
}
