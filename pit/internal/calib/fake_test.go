package calib

import "testing"

func TestRefuseInvented(t *testing.T) {
	h := Card(nil, 30)
	if err := RefuseInvented(h.N, 30, h.Copy); err != nil {
		t.Fatal(err)
	}
	if err := RefuseInvented(0, 30, "Calibration 72%"); err == nil {
		t.Fatal("fake")
	}
}
