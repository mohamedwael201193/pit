package workspace

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/identity"
)

func (s *Store) AssertOwner(wsID string, evm identity.Address) error {
	ws, err := s.Get(wsID)
	if err != nil {
		return err
	}
	if ws.EVM != evm {
		return fmt.Errorf("wrong_user")
	}
	return nil
}
