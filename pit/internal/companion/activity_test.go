package companion

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func TestEvidenceForJobIgnoresStale(t *testing.T) {
	dir := t.TempDir()
	appendActivity(dir, activityEvent{
		Kind: "evidence.filed", Market: "AVAX", JobID: "job-old",
		Root: "0xoldroot", Tx: "0xoldtx", Digest: "0xolddigest",
	})
	appendActivity(dir, activityEvent{
		Kind: "evidence.filed", Market: "HYPE", JobID: "job-live",
		Root: "0xnewroot", Tx: "0xnewtx", Digest: "0xnewdigest", TxLink: "https://chainscan.0g.ai/tx/0xnewtx",
	})
	got := evidenceForJob(dir, "job-live")
	if got == nil || got["tx"] != "0xnewtx" || got["job_id"] != "job-live" {
		t.Fatalf("%v", got)
	}
	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "oldtx") {
		t.Fatal(string(raw))
	}
	if evidenceForJob(dir, "job-missing") != nil {
		t.Fatal("missing job must not fall back")
	}
	if evidenceForJob(dir, "") != nil {
		t.Fatal("empty job")
	}
}

func TestEvidenceForJobDoesNotCrossJobs(t *testing.T) {
	dir := t.TempDir()
	appendActivity(dir, activityEvent{
		Kind: "evidence.filed", Market: "AVAX", JobID: "job-a",
		Root: "0xroota", Tx: "0xtxa", TxLink: "https://chainscan.0g.ai/tx/0xtxa",
	})
	appendActivity(dir, activityEvent{
		Kind: "evidence.filed", Market: "HYPE", JobID: "job-b",
		Root: "0xrootb", Tx: "0xtxb", TxLink: "https://chainscan.0g.ai/tx/0xtxb",
	})
	a := evidenceForJob(dir, "job-a")
	b := evidenceForJob(dir, "job-b")
	if a == nil || a["tx"] != "0xtxa" || a["job_id"] != "job-a" {
		t.Fatalf("a %v", a)
	}
	if b == nil || b["tx"] != "0xtxb" || b["job_id"] != "job-b" {
		t.Fatalf("b %v", b)
	}
	rawA, _ := json.Marshal(a)
	rawB, _ := json.Marshal(b)
	if strings.Contains(string(rawA), "0xtxb") || strings.Contains(string(rawB), "0xtxa") {
		t.Fatal("cross-contaminate")
	}
}
