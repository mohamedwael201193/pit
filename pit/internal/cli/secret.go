package cli

import (
	"fmt"
	"strings"
)

func RefusePrint(s string) error {
	low := strings.ToLower(s)
	for _, w := range []string{"private_key", "mnemonic", "session_key", "hl_secret"} {
		if strings.Contains(low, w) {
			return fmt.Errorf("secret_print")
		}
	}
	return nil
}
