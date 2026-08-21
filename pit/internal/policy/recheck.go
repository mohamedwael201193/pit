package policy

import "fmt"

// Recheck runs after async work. A mutation of clip/side/market fails closed.
func Recheck(before, after Policy, c Context) error {
	if before.Version != after.Version {
		return fmt.Errorf("policy_changed")
	}
	bh, err := before.Hash()
	if err != nil {
		return err
	}
	ah, err := after.Hash()
	if err != nil {
		return err
	}
	if bh != ah {
		return fmt.Errorf("policy_changed")
	}
	return Check(after, c)
}
