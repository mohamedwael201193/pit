package compute

import (
	"fmt"
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
	if _, err := PrepareDirect(j); err != nil {
		return err
	}
	return fmt.Errorf("sealer_exec_not_attached")
}
