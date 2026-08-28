package hl

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) ExtraAgents(user string) (json.RawMessage, error) {
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("user_required")
	}
	return c.postInfo(map[string]any{"type": "extraAgents", "user": user})
}

func SessionAgentLinked(raw json.RawMessage, name, addr string, nowMs int64) bool {
	want := strings.TrimSpace(addr)
	if want == "" {
		return false
	}
	var rows []struct {
		Name       string `json:"name"`
		Address    string `json:"address"`
		ValidUntil *int64 `json:"validUntil"`
	}
	if json.Unmarshal(raw, &rows) != nil {
		return false
	}
	_ = name // Hyperliquid labels are <17 chars; the session address is the authority.
	for _, r := range rows {
		if !strings.EqualFold(r.Address, want) {
			continue
		}
		if r.ValidUntil != nil && *r.ValidUntil > 0 && nowMs >= *r.ValidUntil {
			continue
		}
		return true
	}
	return false
}
