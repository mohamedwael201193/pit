package session

import (
	"fmt"
	"strings"
)

func AgentName(workspaceID string) (string, error) {
	id := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(workspaceID), "-", ""))
	if len(id) < 8 {
		return "", fmt.Errorf("short_workspace")
	}
	return "PIT-" + id[:8], nil
}
