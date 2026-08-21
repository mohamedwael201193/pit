package hl

import (
	"encoding/json"
	"testing"
)

func TestNamedAgentPresent(t *testing.T) {
	raw := json.RawMessage(`[{"name":"PIT-abcd1234","address":"0xabc"},{"name":"other"}]`)
	if !NamedAgentPresent(raw, "PIT-abcd1234") {
		t.Fatal("miss")
	}
	if NamedAgentPresent(raw, "PIT-missing") {
		t.Fatal("ghost")
	}
}
