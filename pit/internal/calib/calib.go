package calib

import "math"

type Sample struct {
	P      float64
	Outcome bool
}

func Brier(samples []Sample) (float64, bool) {
	if len(samples) == 0 {
		return 0, false
	}
	var s float64
	for _, x := range samples {
		y := 0.0
		if x.Outcome {
			y = 1
		}
		d := x.P - y
		s += d * d
	}
	return s / float64(len(samples)), true
}

func ECE(samples []Sample, bins int) (float64, bool) {
	if len(samples) < 1 || bins < 2 {
		return 0, false
	}
	type bin struct {
		n int
		p float64
		y float64
	}
	bs := make([]bin, bins)
	for _, x := range samples {
		i := int(math.Min(float64(bins-1), math.Floor(x.P*float64(bins))))
		bs[i].n++
		bs[i].p += x.P
		if x.Outcome {
			bs[i].y++
		}
	}
	var e float64
	n := float64(len(samples))
	for _, b := range bs {
		if b.n == 0 {
			continue
		}
		avgP := b.p / float64(b.n)
		avgY := b.y / float64(b.n)
		e += (float64(b.n) / n) * math.Abs(avgP-avgY)
	}
	return e, true
}

func Sufficient(n, need int) bool {
	return n >= need
}
