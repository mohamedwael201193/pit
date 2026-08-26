package obs

import "fmt"

func RefuseHealthSecrets(body map[string]any) error {
	for _, k := range []string{"wallet", "session", "private_key", "mnemonic", "session_key", "hl_secret", "authorization", "private_book", "strategy", "app-sk"} {
		if _, ok := body[k]; ok {
			return fmt.Errorf("health_leak")
		}
	}
	return nil
}
