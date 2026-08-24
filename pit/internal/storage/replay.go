package storage

import "fmt"

func RefuseRootReplay(boundRoot, presented string) error {
	if boundRoot == "" || presented == "" || boundRoot != presented {
		return fmt.Errorf("root_replay")
	}
	return nil
}
