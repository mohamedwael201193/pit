package hl

import (
	"encoding/json"
	"testing"
)

func TestSessionAgentLinked(t *testing.T) {
	raw := json.RawMessage(`[{"name":"PIT-abcd1234","address":"0xabc","validUntil":null}]`)
	if !SessionAgentLinked(raw, "PIT-abcd1234", "0xAbC", 1) {
		t.Fatal("miss")
	}
	truncated := json.RawMessage(`[{"name":"-abcd1234","address":"0xabc","validUntil":null}]`)
	if !SessionAgentLinked(truncated, "PIT-abcd1234", "0xabc", 1) {
		t.Fatal("address is enough")
	}
	expired := json.RawMessage(`[{"name":"PIT-abcd1234","address":"0xabc","validUntil":10}]`)
	if SessionAgentLinked(expired, "PIT-abcd1234", "0xabc", 11) {
		t.Fatal("expired")
	}
	if SessionAgentLinked(raw, "PIT-abcd1234", "0xdef", 1) {
		t.Fatal("addr")
	}
}
