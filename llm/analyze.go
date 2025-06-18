// llm/analyze.go
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

func AnalyzePodIssue(prompt string) (string, error) {
	url := "http://localhost:11434/api/generate"
	payload := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(encoded))
	if err != nil {
		return "", fmt.Errorf("post error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	var result OllamaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	return result.Response, nil
}
