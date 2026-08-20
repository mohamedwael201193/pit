package calib

import "testing"

func TestBrierAndEmpty(t *testing.T) {
	if _, ok := Brier(nil); ok {
		t.Fatal("empty")
	}
	v, ok := Brier([]Sample{{P: 1, Outcome: true}, {P: 0, Outcome: false}})
	if !ok || v != 0 {
		t.Fatalf("%v %v", v, ok)
	}
	if Sufficient(10, 30) {
		t.Fatal("must show insufficient data")
	}
}
