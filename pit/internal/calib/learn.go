package calib

func Overconfident(ece, pin float64) bool {
	if pin <= 0 {
		return false
	}
	return ece > pin
}

func Learned(before, after Health) bool {
	if !after.Sufficient {
		return false
	}
	return after.N > before.N
}
