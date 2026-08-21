package keyring

import (
	"fmt"
	"strings"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func Namespace(network, workspaceID, kind string) (string, error) {
	ws, err := identity.ParseWorkspaceID(workspaceID)
	if err != nil {
		return "", err
	}
	if network == "" || kind == "" || strings.Contains(kind, "/") {
		return "", fmt.Errorf("bad_namespace")
	}
	return fmt.Sprintf("pit/%s/%s/%s", network, ws, kind), nil
}
