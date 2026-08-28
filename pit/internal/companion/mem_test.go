package companion

import (
	"os"
	"strings"
	"testing"
)

func TestChatTranscriptIsEncrypted(t *testing.T) {
	dir := t.TempDir()
	appendChat(dir, "user", "secret thesis for A", "")
	raw, err := os.ReadFile(chatPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "secret thesis") {
		t.Fatal("plaintext chat")
	}
	got := readChat(dir, 10)
	if len(got) != 1 || got[0].Text != "secret thesis for A" {
		t.Fatalf("%+v", got)
	}
}

func TestWorkingMemoryHasNoPrompt(t *testing.T) {
	dir := t.TempDir()
	writeWorkingMemory(dir, map[string]any{"coin": "ETH", "kind": "READY_ELIGIBLE"})
	raw, err := os.ReadFile(workingPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "READY_ELIGIBLE") || strings.Contains(strings.ToLower(string(raw)), "prompt") {
		t.Fatal(string(raw))
	}
}
