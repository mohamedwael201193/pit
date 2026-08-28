package compute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
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
	low := strings.ToLower(text + " " + err.Error())
	if strings.Contains(text, "VERIFY_FAIL") {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if strings.Contains(low, "timeout") || strings.Contains(low, "deadline exceeded") || strings.Contains(low, "deadlineexceeded") {
		return fmt.Errorf("DIRECT_PROVIDER_TIMEOUT")
	}
	if ee, ok := err.(*exec.ExitError); ok {
		switch ee.ExitCode() {
		case 10:
			return fmt.Errorf("ROUTER_DOWNGRADE_DENIED")
		case 11:
			return fmt.Errorf("NOT_TEEML")
		case 1:
			if strings.Contains(text, "PUBKEY") || strings.Contains(text, "POST") || strings.Contains(text, "SIG_FAIL") {
				return fmt.Errorf("DIRECT_PROVIDER_TIMEOUT")
			}
			return fmt.Errorf("sealer_runtime")
		case 3:
			if strings.Contains(text, "POST_FAIL 401") || strings.Contains(text, "POST_FAIL 403") || strings.Contains(text, "POST_FAIL 400") || strings.Contains(low, "insufficient balance") {
				return fmt.Errorf("direct_ledger")
			}
			return fmt.Errorf("direct_provider_http")
		case 4:
			return fmt.Errorf("direct_no_chat_id")
		case 5:
			return fmt.Errorf("DIRECT_PROVIDER_TIMEOUT")
		case 6:
			return fmt.Errorf("TEE_VERIFY_FAIL")
		case 7:
			return fmt.Errorf("TEE_OPEN_FAIL")
		}
	}
	if strings.Contains(text, "POST_FAIL 401") || strings.Contains(text, "POST_FAIL 403") || strings.Contains(text, "POST_FAIL 400") || strings.Contains(low, "insufficient balance") {
		return fmt.Errorf("direct_ledger")
	}
	return fmt.Errorf("DIRECT_PROVIDER_UNAVAILABLE")
}

func stopped(stop func() bool) bool {
	return stop != nil && stop()
}

// RunSealedAsk never falls back to Router or plaintext. Missing binary stops the operation.
func RunSealedAsk(j DirectJob) error {
	return RunSealedAskCtl(j, nil, nil)
}

func RunSealedAskCtl(j DirectJob, stage StageFn, stop func() bool) error {
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
	if stopped(stop) {
		return fmt.Errorf("research_cancelled")
	}
	cmd := exec.Command(j.Bin, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Start(); err != nil {
		return sealerExitError(err, buf.Bytes())
	}
	done := make(chan struct{})
	defer close(done)
	if stop != nil {
		go func() {
			tick := time.NewTicker(200 * time.Millisecond)
			defer tick.Stop()
			for {
				select {
				case <-done:
					return
				case <-tick.C:
					if stopped(stop) && cmd.Process != nil {
						_ = cmd.Process.Kill()
						return
					}
				}
			}
		}()
	}
	waitErr := cmd.Wait()
	out := buf.Bytes()
	if stopped(stop) {
		return fmt.Errorf("research_cancelled")
	}
	if waitErr != nil {
		return sealerExitError(waitErr, out)
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
