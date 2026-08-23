package policy

import "fmt"

func MatchPin(pinned, current string) error {
	if pinned == "" || current == "" {
		return fmt.Errorf("pin_incomplete")
	}
	if pinned != current {
		return fmt.Errorf("policy_changed")
	}
	return nil
}
