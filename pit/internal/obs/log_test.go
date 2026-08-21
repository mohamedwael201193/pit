package obs

import (
	"bytes"
	"strings"
	"testing"
)

func TestSanitizeSecrets(t *testing.T) {
	if Sanitize("session_key leaked") != "redacted" {
		t.Fatal("session")
	}
	if Sanitize("ok") != "ok" {
		t.Fatal("ok")
	}
	var buf bytes.Buffer
	if err := Write(&buf, Event{RequestID: "r1", Phase: "POLICY_CHECK", Status: "ok"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"requestId":"r1"`) {
		t.Fatal(buf.String())
	}
}
