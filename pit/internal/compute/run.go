package compute

import (
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
	if strings.HasSuffix(low, ".ts") || strings.HasSuffix(low, ".js") {
		return fmt.Errorf("script_sealer_denied")
	}
	return nil
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
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if err := RequireScheme(string(out)); err != nil {
		return err
	}
	return nil
}
