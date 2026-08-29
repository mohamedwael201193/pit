package hl

import (
	"bytes"
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

func (c *Client) UserFills(user string) (json.RawMessage, error) {
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("user_required")
	}
	return c.postInfo(map[string]any{"type": "userFills", "user": user})
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

func decodeRows(raw json.RawMessage) []map[string]any {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var rows []map[string]any
	if dec.Decode(&rows) != nil {
		return nil
	}
	return rows
}

func oidField(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case json.Number:
		return strings.TrimSpace(t.String())
	case string:
		return strings.Trim(strings.TrimSpace(t), `"`)
	default:
		s := strings.TrimSpace(fmt.Sprint(t))
		s = strings.Trim(s, `"`)
		if s == "<nil>" {
			return ""
		}
		return s
	}
}

func OIDOnVenue(raw json.RawMessage, oid string) bool {
	want := strings.TrimSpace(oid)
	if want == "" {
		return false
	}
	for _, r := range decodeRows(raw) {
		if oidField(r["oid"]) == want {
			return true
		}
	}
	return false
}

func OIDInFills(raw json.RawMessage, oid string) bool {
	want := strings.TrimSpace(oid)
	if want == "" {
		return false
	}
	for _, r := range decodeRows(raw) {
		if oidField(r["oid"]) == want {
			return true
		}
	}
	return false
}
