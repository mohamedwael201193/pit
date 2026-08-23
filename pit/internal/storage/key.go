package storage

import "fmt"

func KeysMatch(want, got string) error {
	if err := RequireHexKey(want); err != nil {
		return err
	}
	if err := RequireHexKey(got); err != nil {
		return err
	}
	if want != got {
		return fmt.Errorf("wrong_key")
	}
	return nil
}
