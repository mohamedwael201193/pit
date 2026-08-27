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

func sealerExitError(err error, out []byte) error {
	if err == nil {
		return nil
	}
	text := string(out)
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {
		case 10:
			return fmt.Errorf("ROUTER_DOWNGRADE_DENIED")
		case 11:
			return fmt.Errorf("NOT_TEEML")
		case 1:
			return fmt.Errorf("sealer_runtime")
		case 3:
			if strings.Contains(text, "POST_FAIL 401") || strings.Contains(text, "POST_FAIL 403") {
				return fmt.Errorf("direct_ledger")
			}
			return fmt.Errorf("direct_provider_http")
		case 4:
			return fmt.Errorf("direct_no_chat_id")
		case 5:
			return fmt.Errorf("direct_signature_http")
		case 6:
			return fmt.Errorf("TEE_VERIFY_FAIL")
		case 7:
			return fmt.Errorf("TEE_OPEN_FAIL")
		}
	}
	if strings.Contains(text, "POST_FAIL 401") || strings.Contains(text, "POST_FAIL 403") {
		return fmt.Errorf("direct_ledger")
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
		return sealerExitError(err, out)
	}
	if err := RequireScheme(string(out)); err != nil {
		return err
	}
	return AcceptSealedEvidence(j.OutPath, j.OnchainSigner)
}

func PublicRoleEvidence(j DirectJob) map[string]any {
	b, err := os.ReadFile(j.OutPath)
	if err != nil {
		return map[string]any{"role": string(j.Role), "verify_e2ee": "FAIL"}
	}
	var ev sealedEvidence
	if err := json.Unmarshal(b, &ev); err != nil {
		return map[string]any{"role": string(j.Role), "verify_e2ee": "FAIL"}
	}
	return map[string]any{
		"role":          string(j.Role),
		"verify_e2ee":   ev.VerifyE2EE,
		"pubkey_signer": ev.PubkeySigner,
		"teeSigner":     ev.TeeSigner,
	}
}
