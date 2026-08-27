package compute

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func MustNativeSealer(bin string) error {
	b := strings.TrimSpace(bin)
	if b == "" {
		return fmt.Errorf("sealer_not_wired")
	}
	low := strings.ToLower(b)
	if strings.HasSuffix(low, ".py") || strings.Contains(low, "python") {
		return fmt.Errorf("python_gate_not_product")
	}
	if strings.HasSuffix(low, ".ts") || strings.HasSuffix(low, ".js") || strings.HasSuffix(low, ".mjs") || strings.HasSuffix(low, ".cjs") {
		return fmt.Errorf("script_sealer_denied")
	}
	return nil
}

type sealedEvidence struct {
	VerifyE2EE   string `json:"verify_e2ee"`
	SigText      string `json:"sig_text"`
	PubkeySigner string `json:"pubkey_signer"`
	TeeSigner    string `json:"teeSigner"`
}

func AcceptSealedEvidence(outPath, onchain string) error {
	b, err := os.ReadFile(outPath)
	if err != nil {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	var ev sealedEvidence
	if err := json.Unmarshal(b, &ev); err != nil {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if !strings.EqualFold(strings.TrimSpace(ev.VerifyE2EE), "OK") {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if err := RequireScheme(ev.SigText); err != nil {
		return err
	}
	if err := RequireSigner(ev.PubkeySigner, onchain); err != nil {
		return err
	}
	if err := RequireSigner(ev.TeeSigner, onchain); err != nil {
		return err
	}
	return nil
}

func sealerExitError(err error) error {
	if err == nil {
		return nil
	}
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {
		case 10:
			return fmt.Errorf("ROUTER_DOWNGRADE_DENIED")
		case 11:
			return fmt.Errorf("NOT_TEEML")
		}
	}
	return fmt.Errorf("TEE_VERIFY_FAIL")
}

// RunSealedAsk never falls back to Router or plaintext. Missing binary stops the operation.
func RunSealedAsk(j DirectJob) error {
	if err := MustNativeSealer(j.Bin); err != nil {
		return err
	}
	args, err := PrepareDirect(j)
	if err != nil {
		return err
	}
	if _, err := os.Stat(j.Bin); err != nil {
		return fmt.Errorf("sealer_not_wired")
	}
	cmd := exec.Command(j.Bin, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return sealerExitError(err)
	}
	if err := RequireScheme(string(out)); err != nil {
		return err
	}
	return AcceptSealedEvidence(j.OutPath, j.OnchainSigner)
}
