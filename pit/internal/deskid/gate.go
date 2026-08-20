package deskid

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
)

func TransferEnabled(net config.Network) bool {
	return false // Aristotle attestor absent. Galileo path is not claimed live without a PIT tx.
}

func MustAuthorized(owner, caller string, allowlist []string) error {
	if strings.EqualFold(owner, caller) {
		return nil
	}
	for _, a := range allowlist {
		if strings.EqualFold(a, caller) {
			return nil
		}
	}
	return fmt.Errorf("not_authorized")
}

func RefuseTransfer(net config.Network) error {
	if TransferEnabled(net) {
		return nil
	}
	if net == config.Mainnet {
		return fmt.Errorf("AttestorNotOnAristotle")
	}
	return fmt.Errorf("transfer_unverified")
}
