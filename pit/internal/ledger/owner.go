package ledger

import "fmt"

func (s *Store) RefuseForeign(workspace string) error {
	if s == nil || s.workspace == "" {
		return fmt.Errorf("workspace required")
	}
	if workspace != s.workspace {
		return fmt.Errorf("wrong_workspace")
	}
	return nil
}
