package session

import (
	"fmt"
	"strings"
)

func SwitchAccount(sessionWallet, connectedWallet string) error {
	a := strings.ToLower(strings.TrimSpace(sessionWallet))
	b := strings.ToLower(strings.TrimSpace(connectedWallet))
	if a == "" || b == "" {
		return fmt.Errorf("unbound_wallet")
	}
	if a != b {
		return fmt.Errorf("account_switch_requires_new_session")
	}
	return nil
}
