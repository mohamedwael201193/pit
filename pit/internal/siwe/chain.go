package siwe

import "fmt"

func ChainMatch(got, want int64) error {
	if want == 0 || got == 0 {
		return fmt.Errorf("chain_unbound")
	}
	if got != want {
		return fmt.Errorf("WRONG_NETWORK")
	}
	return nil
}
