package exec

import (
	"testing"

	"github.com/mohamedwael201193/pit/internal/policy"
)

func TestRefuseCalibrationFloor(t *testing.T) {
	p := policy.Default()
	p.MinSkillCalibration = 0.65
	if err := RefuseCalibrationFloor(p, 0.2); err == nil {
		t.Fatal("floor")
	}
	if err := RefuseCalibrationFloor(p, 0.7); err != nil {
		t.Fatal(err)
	}
}
