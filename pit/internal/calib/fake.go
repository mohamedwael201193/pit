package calib

import (
	"fmt"
	"strings"
)

func RefuseInvented(n, need int, copy string) error {
	if n >= need {
		return nil
	}
	low := strings.ToLower(copy)
	if strings.Contains(copy, "72%") || strings.Contains(copy, "99.99%") {
		return fmt.Errorf("fake_calibration")
	}
	if !strings.Contains(low, "not enough") {
		return fmt.Errorf("empty_copy")
	}
	return nil
}
