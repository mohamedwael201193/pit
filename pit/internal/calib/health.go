package calib

import "fmt"

type Health struct {
	N          int
	Brier      float64
	ECE        float64
	Sufficient bool
	Copy       string
}

func Card(samples []Sample, need int) Health {
	n := len(samples)
	h := Health{N: n, Sufficient: Sufficient(n, need)}
	if !h.Sufficient {
		h.Copy = fmt.Sprintf("Not enough samples (%d/%d). PIT will not invent accuracy.", n, need)
		return h
	}
	b, _ := Brier(samples)
	e, _ := ECE(samples, 10)
	h.Brier, h.ECE = b, e
	h.Copy = fmt.Sprintf("Brier %.3f · ECE %.3f · N=%d", b, e, n)
	if e > 0.15 {
		h.Copy += " · overconfidence"
	}
	return h
}
