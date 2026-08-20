package watch

import "fmt"

func CloudLoopForbidden() error {
	return fmt.Errorf("no_cloud_watch")
}

func GhostCard(nReal int) error {
	if nReal < 0 {
		return fmt.Errorf("ghost_card")
	}
	return nil
}
