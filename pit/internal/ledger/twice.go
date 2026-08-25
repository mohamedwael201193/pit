package ledger

import "fmt"

func RefuseSecondApply(applied bool) error {
	if !applied {
		return fmt.Errorf("duplicate_action")
	}
	return nil
}
