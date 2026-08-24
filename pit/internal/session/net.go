package session

import (
	"fmt"
	"strings"
)

func MatchNetwork(sessionNet, workspaceNet string) error {
	a := strings.ToLower(strings.TrimSpace(sessionNet))
	b := strings.ToLower(strings.TrimSpace(workspaceNet))
	if a == "" || b == "" {
		return fmt.Errorf("network_required")
	}
	if a != b {
		return fmt.Errorf("session_network_mismatch")
	}
	return nil
}
