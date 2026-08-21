package hl

import (
	"encoding/json"
	"strings"
)

func NamedAgentPresent(raw json.RawMessage, want string) bool {
	var rows []struct {
		Name  string `json:"name"`
		Valid bool   `json:"valid,omitempty"`
	}
	if json.Unmarshal(raw, &rows) != nil {
		return false
	}
	for _, r := range rows {
		if strings.EqualFold(r.Name, want) {
			return true
		}
	}
	return false
}
