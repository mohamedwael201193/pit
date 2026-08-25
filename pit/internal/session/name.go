package session

import (
	"fmt"
	"strings"
)

func RequireAgentName(name, workspaceID string) error {
	want := AgentKey{}.Name(workspaceID)
	if !strings.HasPrefix(name, "PIT-") || name != want {
		return fmt.Errorf("agent_name")
	}
	return nil
}
