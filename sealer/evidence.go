package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
)

func sha256Hex(b []byte) string {
	s := sha256.Sum256(b)
	return hex.EncodeToString(s[:])
}

func writeEvidence(path string, result map[string]any) error {
	if _, ok := result["prompt"]; ok {
		delete(result, "prompt")
	}
	if _, ok := result["authorization"]; ok {
		delete(result, "authorization")
	}
	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o600)
}

func extractText(opened map[string]json.RawMessage) string {
	raw, ok := opened["choices"]
	if !ok {
		return ""
	}
	var choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.Unmarshal(raw, &choices); err != nil || len(choices) == 0 {
		return strings.TrimSpace(string(raw))
	}
	return strings.TrimSpace(choices[0].Message.Content)
}

func asEnv(raw []byte) (map[string]json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}
