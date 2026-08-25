package wallet

import "fmt"

func RefuseDeclinedBind(accepted bool) error {
	if !accepted {
		return fmt.Errorf("siwe_declined")
	}
	return nil
}
