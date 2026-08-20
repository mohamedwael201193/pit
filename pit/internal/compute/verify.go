package compute

import (
	"fmt"
	"strings"
)

const SchemeE2EE = "zg-sig-v1/e2ee-ct"

func RequireScheme(text string) error {
	if !strings.HasPrefix(text, SchemeE2EE) {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	return nil
}

func RequireSigner(recovered, onchain string) error {
	if recovered == "" || onchain == "" {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	if !strings.EqualFold(strings.TrimSpace(recovered), strings.TrimSpace(onchain)) {
		return fmt.Errorf("TEE_VERIFY_FAIL")
	}
	return nil
}

func RejectPlaintextFallback(allowFallbacks string) error {
	v := strings.ToLower(strings.TrimSpace(allowFallbacks))
	if v != "" && v != "false" && v != "0" {
		return fmt.Errorf("plaintext_fallback_denied")
	}
	return nil
}
