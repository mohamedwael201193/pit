package compute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMustNativeSealer(t *testing.T) {
	if err := MustNativeSealer(""); err == nil {
		t.Fatal("empty")
	}
	if err := MustNativeSealer("committee.py"); err == nil {
		t.Fatal("python")
	}
	if err := MustNativeSealer("sealer.ts"); err == nil {
		t.Fatal("ts")
	}
	if err := MustNativeSealer("/usr/local/bin/pit-sealer"); err != nil {
		t.Fatal(err)
	}
}

func TestRunSealedAskMissingBinary(t *testing.T) {
	err := RunSealedAsk(DirectJob{
		Bin:           filepath.Join(t.TempDir(), "no-such-sealer"),
		AuthPath:      filepath.Join(t.TempDir(), "a"),
		PromptPath:    filepath.Join(t.TempDir(), "p"),
		OutPath:       filepath.Join(t.TempDir(), "o"),
		Role:          Researcher,
		ProviderURL:   "https://compute-network-19.integratenetwork.work",
		OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err == nil || err.Error() != "sealer_not_wired" {
		t.Fatalf("%v", err)
	}
}

func TestRunSealedAskRejectsPlainOutput(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sealer.bat")
	if err := os.WriteFile(bin, []byte("echo plaintext\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	err := RunSealedAsk(DirectJob{
		Bin:           bin,
		AuthPath:      filepath.Join(t.TempDir(), "a"),
		PromptPath:    filepath.Join(t.TempDir(), "p"),
		OutPath:       filepath.Join(t.TempDir(), "o"),
		Role:          Researcher,
		ProviderURL:   "https://compute-network-19.integratenetwork.work",
		OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err == nil {
		t.Fatal("plain")
	}
}
