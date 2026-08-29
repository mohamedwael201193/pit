package companion

import (
	"os"
	"path/filepath"
	"testing"
)

func TestActivityStoresLink(t *testing.T) {
	dir := t.TempDir()
	appendActivity(dir, activityEvent{Kind: "opportunity", Market: "BTC", Link: "https://app.hyperliquid.xyz/trade/BTC"})
	got := readActivity(dir, 10)
	if len(got) != 1 || got[0].Link != "https://app.hyperliquid.xyz/trade/BTC" {
		t.Fatalf("%+v", got)
	}
}

func TestActivityDoesNotStoreSecrets(t *testing.T) {
	dir := t.TempDir()
	appendActivity(dir, activityEvent{Kind: "research.started", Market: "ETH", JobID: "j1", Action: "research"})
	got := readActivity(dir, 10)
	if len(got) != 1 || got[0].Kind != "research.started" {
		t.Fatalf("%+v", got)
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "activity.jsonl"))
	if string(raw) == "" {
		t.Fatal("empty")
	}
}

func TestActivityDropsAppSk(t *testing.T) {
	dir := t.TempDir()
	appendActivity(dir, activityEvent{Kind: "leak", Reason: "Bearer app-sk-secret"})
	if len(readActivity(dir, 10)) != 0 {
		t.Fatal("stored secret")
	}
}
