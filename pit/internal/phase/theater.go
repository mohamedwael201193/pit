package phase

import "fmt"

func RefuseTheater(s string) error {
	switch s {
	case "SPINNING", "LOADING_FAKE", "setTimeout":
		return fmt.Errorf("theater")
	}
	if !Known(s) {
		return fmt.Errorf("unknown_phase")
	}
	return nil
}
