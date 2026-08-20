package calib

import (
	"strings"
	"testing"
)

func TestHealthEmptyCopy(t *testing.T) {
	h := Card(nil, 30)
	if h.Sufficient || h.Copy == "" {
		t.Fatal(h)
	}
	if h.Brier != 0 {
		t.Fatal("must not invent accuracy")
	}
}

func TestOverconfidenceFlag(t *testing.T) {
	samples := make([]Sample, 30)
	for i := range samples {
		samples[i] = Sample{P: 0.95, Outcome: false}
	}
	h := Card(samples, 30)
	if !h.Sufficient || !strings.Contains(h.Copy, "overconfidence") {
		t.Fatal(h)
	}
}
