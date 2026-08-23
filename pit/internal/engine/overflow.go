package engine

import "fmt"

func ClipNotExceed(notional, maxClip float64) error {
	if err := RejectNonFinite(notional, "notional"); err != nil {
		return err
	}
	if err := RejectNonFinite(maxClip, "clip"); err != nil {
		return err
	}
	if maxClip <= 0 {
		return fmt.Errorf("bad_clip")
	}
	if notional > maxClip+1e-9 {
		return fmt.Errorf("notional_exceeds_clip")
	}
	return nil
}
