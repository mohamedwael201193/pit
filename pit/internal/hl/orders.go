package hl

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) OpenOrders(user string) (json.RawMessage, error) {
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("user_required")
	}
	return c.postInfo(map[string]any{"type": "frontendOpenOrders", "user": user})
}

func CloidOnVenue(raw json.RawMessage, cloid string) bool {
	var rows []struct {
		Cloid string `json:"cloid"`
		C     string `json:"c"`
	}
	if json.Unmarshal(raw, &rows) != nil {
		return false
	}
	want := strings.ToLower(strings.TrimSpace(cloid))
	for _, r := range rows {
		got := strings.ToLower(strings.TrimSpace(r.Cloid))
		if got == "" {
			got = strings.ToLower(strings.TrimSpace(r.C))
		}
		if got != "" && got == want {
			return true
		}
	}
	return false
}
