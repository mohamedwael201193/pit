package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func RefuseCalibrationFloor(p policy.Policy, skill float64) error {
	if p.MinSkillCalibration > 0 && skill+1e-12 < p.MinSkillCalibration {
		return fmt.Errorf("calibration_floor")
	}
	return nil
}
