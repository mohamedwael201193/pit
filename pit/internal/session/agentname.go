package session

import (
	"fmt"
	"strings"
)

const MaxAgentNameLen = 16 // Hyperliquid: API wallet name must be less than 17 characters.

func AgentName(workspaceID string) (string, error) {
	id := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(workspaceID), "-", ""))
	if len(id) < 8 {
		return "", fmt.Errorf("short_workspace")
	}
	name := "PIT-" + id[:8]
	if len(name) > MaxAgentNameLen {
		return "", fmt.Errorf("agent_name")
	}
	return name, nil
}
