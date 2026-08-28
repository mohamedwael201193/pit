package engine

import (
	"encoding/json"
	"strings"
)

// ParseRoleText extracts a host-readable role object from model text.
// Size, leverage, and permissions in the payload are ignored by BindPreview.
func ParseRoleText(text string) RoleJSON {
	s := strings.TrimSpace(text)
	if s == "" {
		return RoleJSON{}
	}
	if i := strings.Index(s, "```"); i >= 0 {
		rest := strings.TrimSpace(s[i+3:])
		rest = strings.TrimPrefix(rest, "json")
		rest = strings.TrimSpace(rest)
		if j := strings.Index(rest, "```"); j >= 0 {
			s = rest[:j]
		} else {
			s = rest
		}
	}
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start < 0 || end <= start {
		return RoleJSON{}
	}
	var r RoleJSON
	_ = json.Unmarshal([]byte(s[start:end+1]), &r)
	return r
}
