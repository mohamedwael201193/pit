package storage

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/config"
	"github.com/mohamedwael201193/pit/internal/identity"
)

func ObjectKey(net config.Network, workspaceID, kind, name string) (string, error) {
	ws, err := identity.ParseWorkspaceID(workspaceID)
	if err != nil {
		return "", err
	}
	if kind == "" || strings.Contains(kind, "/") || strings.Contains(name, "/") {
		return "", fmt.Errorf("bad object name")
	}
	return fmt.Sprintf("%s/ws/%s/%s/%s", net, ws, kind, name), nil
}

func AssertWorkspace(key, workspaceID string) error {
	ws, err := identity.ParseWorkspaceID(workspaceID)
	if err != nil {
		return err
	}
	needle := "/ws/" + ws + "/"
	if !strings.Contains(key, needle) {
		return fmt.Errorf("wrong_workspace")
	}
	return nil
}

func RequireHexKey(key string) error {
	if !strings.HasPrefix(key, "0x") || len(key) != 66 {
		return fmt.Errorf("encryption key must be 0x-prefixed 32 bytes")
	}
	return nil
}
