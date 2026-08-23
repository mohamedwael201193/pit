package watch

import "github.com/mohamedwael201193/pit/internal/identity"

func Bound(workspaceID string) error {
	_, err := identity.ParseWorkspaceID(workspaceID)
	return err
}
