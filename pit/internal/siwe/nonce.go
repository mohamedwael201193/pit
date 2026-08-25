package siwe

import "fmt"

func RefuseNonceReplay(used map[string]struct{}, nonce string) error {
	if nonce == "" {
		return fmt.Errorf("nonce_required")
	}
	if _, ok := used[nonce]; ok {
		return fmt.Errorf("nonce_replay")
	}
	return nil
}
