package compute

import "testing"

func TestRequireSealedMessages(t *testing.T) {
	if err := RequireSealedMessages([]string{"tools"}); err == nil {
		t.Fatal("messages")
	}
	if err := RequireSealedMessages([]string{"messages", "tools"}); err != nil {
		t.Fatal(err)
	}
	if err := RefuseCleartextMessages(map[string]any{"messages": []any{}}); err == nil {
		t.Fatal("plain")
	}
	if err := RefuseCleartextMessages(map[string]any{"messages": []any{}, E2EEKey: map[string]any{}}); err != nil {
		t.Fatal(err)
	}
}
