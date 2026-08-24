package engine

import "fmt"

func RefuseUncertainty(u, max float64) error {
	if err := RejectNonFinite(u, "uncertainty"); err != nil {
		return err
	}
	if max > 0 && u > max {
		return fmt.Errorf("uncertainty")
	}
	return nil
}
