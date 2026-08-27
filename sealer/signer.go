package main

import (
	"fmt"
	"strings"
)

func requirePubKeyMatchesOnchain(pubSigner, onchain string) error {
	pub := strings.TrimSpace(pubSigner)
	tee := strings.TrimSpace(onchain)
	if pub == "" || tee == "" {
		return fmt.Errorf("tee_signer_mismatch")
	}
	if !strings.EqualFold(pub, tee) {
		return fmt.Errorf("tee_signer_mismatch")
	}
	return nil
}
