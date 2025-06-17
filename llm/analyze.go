package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaResponse struct {
	Response string `json:"response"`
}

// AnalyzePodIssue sends a prompt to the local Ollama LLM and returns the AI's suggestion.
func AnalyzePodIssue(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	payload := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal prompt: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama LLM: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read LLM response: %w", err)
	}

	var result OllamaResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal LLM response: %w", err)
	}

	return result.Response, nil
}