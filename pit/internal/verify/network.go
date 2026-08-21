package verify

import (
	"fmt"
	"strings"
)

func SameNetwork(receiptNetwork, workspaceNetwork string) error {
	a := strings.ToLower(strings.TrimSpace(receiptNetwork))
	b := strings.ToLower(strings.TrimSpace(workspaceNetwork))
	if a == "" || b == "" {
		return fmt.Errorf("receipt_unbound")
	}
	if a != b {
		return fmt.Errorf("network_mix")
	}
	return nil
}
