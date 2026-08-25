package ledger

import "fmt"

func RefuseEmptyWorkspace(workspace string) error {
	if workspace == "" {
		return fmt.Errorf("workspace required")
	}
	return nil
}
