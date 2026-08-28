package compute

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestAcceptSealedEvidence(t *testing.T) {
	p := filepath.Join(t.TempDir(), "ev.json")
	body := `{"verify_e2ee":"OK","sig_text":"zg-sig-v1/e2ee-ct:aa","pubkey_signer":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9","teeSigner":"0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9"}`
	if err := os.WriteFile(p, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := AcceptSealedEvidence(p, "0xa46ea4fc5889ad35a1487e1ed04dccfa872146b9"); err != nil {
		t.Fatal(err)
	}
	if err := AcceptSealedEvidence(p, "0x0000000000000000000000000000000000000001"); err == nil {
		t.Fatal("wrong onchain")
	}
}

func TestSealerExitErrorLedger(t *testing.T) {
	err := sealerExitError(fmt.Errorf("x"), []byte("POST_FAIL 401\n"))
	if err == nil || err.Error() != "direct_ledger" {
		t.Fatalf("%v", err)
	}
	err = sealerExitError(fmt.Errorf("x"), []byte("POST_FAIL 400\n"))
	if err == nil || err.Error() != "direct_ledger" {
		t.Fatalf("400 %v", err)
	}
	err = sealerExitError(fmt.Errorf("x"), []byte("VERIFY_FAIL signer"))
	if err == nil || err.Error() != "TEE_VERIFY_FAIL" {
		t.Fatalf("%v", err)
	}
}

func TestRequireSchemeRejectsHumanStdout(t *testing.T) {
	if err := RequireScheme("COMMITTEE_OK researcher zg-sig-v1/e2ee-ct:aa"); err == nil {
		t.Fatal("human stdout is not the signature")
	}
}

func TestRunSealedAskHonorsEvidenceNotStdoutPrefix(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "ok-sealer")
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, "./testdata/oksealer")
	cmd.Dir = "."
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build fake sealer: %v %s", err, out)
	}
	outPath := filepath.Join(dir, "researcher.json")
	err := RunSealedAsk(DirectJob{
		Bin:           bin,
		AuthPath:      filepath.Join(dir, "a"),
		PromptPath:    filepath.Join(dir, "p"),
		OutPath:       outPath,
		Role:          Researcher,
		ProviderURL:   "https://compute-network-19.integratenetwork.work",
		OnchainSigner: "0xA46EA4FC5889AD35A1487e1Ed04dCcfa872146B9",
	})
	if err != nil {
		t.Fatal(err)
	}
}

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
	if err := MustNativeSealer("sealer.mjs"); err == nil {
		t.Fatal("mjs")
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
