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

func SessionAgentUntil(raw json.RawMessage, name, addr string, nowMs int64) (bool, int64) {
	want := strings.TrimSpace(addr)
	if want == "" {
		return false, 0
	}
	var rows []struct {
		Name       string `json:"name"`
		Address    string `json:"address"`
		ValidUntil *int64 `json:"validUntil"`
	}
	if json.Unmarshal(raw, &rows) != nil {
		return false, 0
	}
	_ = name // Hyperliquid labels are <17 chars; the session address is the authority.
	for _, r := range rows {
		if !strings.EqualFold(r.Address, want) {
			continue
		}
		until := int64(0)
		if r.ValidUntil != nil {
			until = *r.ValidUntil
		}
		if until > 0 && nowMs >= until {
			continue
		}
		return true, until
	}
	return false, 0
}

func SessionAgentLinked(raw json.RawMessage, name, addr string, nowMs int64) bool {
	ok, _ := SessionAgentUntil(raw, name, addr, nowMs)
	return ok
}
