package storage

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TestLiveOfficialUploadDownloadProof(t *testing.T) {
	if os.Getenv("PIT_LIVE_STORAGE") != "1" {
		t.Skip("set PIT_LIVE_STORAGE=1 to run against Aristotle turbo storage")
	}
	cli := LookCLI()
	if err := RefuseMissingProof(cli); err != nil {
		t.Fatal(err)
	}
	payer := strings.TrimSpace(os.Getenv("PIT_OG_PAYER_KEY"))
	norm, err := NormalizePayerKey(payer)
	if err != nil {
		t.Fatal("payer_key_required")
	}
	payer = norm
	var raw [32]byte
	if _, err := rand.Read(raw[:]); err != nil {
		t.Fatal(err)
	}
	mem := "0x" + hex.EncodeToString(raw[:])
	if mem == payer {
		t.Fatal("generated key collided with payer")
	}
	dir := t.TempDir()
	in := filepath.Join(dir, "payload.bin")
	out := filepath.Join(dir, "out.bin")
	wrong := filepath.Join(dir, "wrong.bin")
	payload := []byte("pit-live-storage-proof-v0.1.0")
	if err := os.WriteFile(in, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	ch := config.For(config.Mainnet)
	upJob := Job{
		CLI: cli, RPC: ch.RPC, Indexer: ch.StorageIndexer, Flow: ch.StorageFlow,
		KeyHex: mem, PayerKey: payer, InputPath: in,
	}
	upArgs, err := UploadArgs(upJob)
	if err != nil {
		t.Fatal(err)
	}
	cmd := Command(upJob, upArgs)
	cmd.Dir = dir
	outb, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("upload failed (output redacted length %d): %v", len(outb), err)
	}
	root := extractRoot(outb)
	if err := RejectBadRoot(root); err != nil {
		t.Fatalf("no root in upload output")
	}
	downJob := Job{CLI: cli, Indexer: ch.StorageIndexer, KeyHex: mem, Root: root, OutPath: out}
	downArgs, err := DownloadArgs(downJob)
	if err != nil {
		t.Fatal(err)
	}
	dcmd := Command(downJob, downArgs)
	dcmd.Dir = dir
	if _, err := dcmd.CombinedOutput(); err != nil {
		t.Fatal("download --proof failed")
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, payload) {
		t.Fatal("roundtrip mismatch")
	}
	bad := Job{CLI: cli, Indexer: ch.StorageIndexer, KeyHex: "0x" + strings.Repeat("cd", 32), Root: root, OutPath: wrong}
	bargs, err := DownloadArgs(bad)
	if err != nil {
		t.Fatal(err)
	}
	bcmd := exec.Command(cli, bargs...)
	bcmd.Dir = dir
	_ = bcmd.Run()
	bgot, _ := os.ReadFile(wrong)
	if bytes.Equal(bgot, payload) {
		t.Fatal("wrong encryption key produced plaintext")
	}
	tamperPath := filepath.Join(dir, "tamper.bin")
	tamper := Job{CLI: cli, Indexer: ch.StorageIndexer, KeyHex: mem, Root: "0x" + strings.Repeat("ee", 32), OutPath: tamperPath}
	targs, err := DownloadArgs(tamper)
	if err != nil {
		t.Fatal(err)
	}
	tcmd := exec.Command(cli, targs...)
	tcmd.Dir = dir
	terr := tcmd.Run()
	tgot, _ := os.ReadFile(tamperPath)
	if terr == nil && bytes.Equal(tgot, payload) {
		t.Fatal("wrong root produced original payload")
	}
}

func extractRoot(out []byte) string {
	re := regexp.MustCompile(`0x[0-9a-fA-F]{64}`)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(strings.ToLower(line), "root") {
			if m := re.FindString(line); m != "" {
				return m
			}
		}
	}
	return re.FindString(string(out))
}
