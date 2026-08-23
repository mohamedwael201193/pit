package cli

import (
	"github.com/mohamedwael201193/pit/internal/identity"
	"github.com/mohamedwael201193/pit/internal/policy"
)

func PinWorkspace(dir, workspaceID string, p policy.Policy) (string, error) {
	ws, err := identity.ParseWorkspaceID(workspaceID)
	if err != nil {
		return "", err
	}
	h, err := p.Hash()
	if err != nil {
		return "", err
	}
	return policy.PinFile(dir, ws, h)
}

func CheckPinned(dir, workspaceID string, p policy.Policy) error {
	ws, err := identity.ParseWorkspaceID(workspaceID)
	if err != nil {
		return err
	}
	got, err := policy.ReadPin(dir, ws)
	if err != nil {
		return err
	}
	h, err := p.Hash()
	if err != nil {
		return err
	}
	return policy.MatchPin(got, h)
}
