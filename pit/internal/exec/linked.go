package exec

import "fmt"

func RefusePostUntilLinked(linked bool) error {
	if !linked {
		return fmt.Errorf("approveAgent_required")
	}
	return nil
}
