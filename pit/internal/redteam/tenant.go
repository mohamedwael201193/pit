package redteam

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func TenantMix(a, b string) error {
	wa, err := identity.ParseWorkspaceID(a)
	if err != nil {
		return err
	}
	wb, err := identity.ParseWorkspaceID(b)
	if err != nil {
		return err
	}
	if wa == wb {
		return fmt.Errorf("tenant_collision")
	}
	return nil
}
