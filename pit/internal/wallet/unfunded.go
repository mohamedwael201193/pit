package wallet

import "fmt"

func MapUnfunded(funded bool) error {
	if !funded {
		return fmt.Errorf(HLUnfunded)
	}
	return nil
}
