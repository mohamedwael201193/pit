package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteEvidenceStripsPromptAndAuthorization(t *testing.T) {
	p := filepath.Join(t.TempDir(), "out.json")
	if err := writeEvidence(p, map[string]any{
		"prompt":        "PRIVATE BOOK",
		"authorization": "Bearer app-sk-secret",
		"prompt_sha256": "abc",
		"role":          "researcher",
	}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if _, ok := m["prompt"]; ok {
		t.Fatal("prompt leaked")
	}
	if _, ok := m["authorization"]; ok {
		t.Fatal("authorization leaked")
	}
	if m["prompt_sha256"] != "abc" {
		t.Fatalf("%v", m)
	}
}

func TestSha256HexStable(t *testing.T) {
	if sha256Hex([]byte("PIT")) == "" {
		t.Fatal("hash")
	}
	if sha256Hex([]byte("PIT")) != sha256Hex([]byte("PIT")) {
		t.Fatal("stable")
	}
}
