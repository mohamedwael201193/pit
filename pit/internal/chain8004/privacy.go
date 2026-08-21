package chain8004

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Feedback is public calibration. It must not contain the private book, policy, or strategy.
type Feedback struct {
	AgentID string   `json:"agentId"`
	Score   int      `json:"score"`
	Tags    []string `json:"tags,omitempty"`
}

var banned = []string{"book", "strategy", "policy", "memory", "position", "preview"}

func EncodeFeedback(f Feedback, extra map[string]any) ([]byte, error) {
	if f.AgentID == "" {
		return nil, fmt.Errorf("agent_required")
	}
	if f.Score < 0 || f.Score > 100 {
		return nil, fmt.Errorf("score_range")
	}
	for k := range extra {
		low := strings.ToLower(k)
		for _, b := range banned {
			if strings.Contains(low, b) {
				return nil, fmt.Errorf("private_field_forbidden")
			}
		}
	}
	body := map[string]any{"agentId": f.AgentID, "score": f.Score, "tags": f.Tags}
	for k, v := range extra {
		body[k] = v
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	low := strings.ToLower(string(encoded))
	for _, b := range banned {
		if strings.Contains(low, `"`+b+`"`) {
			return nil, fmt.Errorf("private_field_forbidden")
		}
	}
	return encoded, nil
}
