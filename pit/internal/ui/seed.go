package ui

import (
	"fmt"
	"strings"
)

func RefuseSeedPrompt(s string) error {
	low := strings.ToLower(s)
	for _, w := range []string{"seed phrase", "mnemonic", "private key"} {
		if strings.Contains(low, w) && strings.Contains(low, "enter") {
			return fmt.Errorf("seed_prompt")
		}
	}
	return nil
}
