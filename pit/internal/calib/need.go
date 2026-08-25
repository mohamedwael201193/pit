package calib

func NeedResolved() int { return 30 }

func RefuseSparse(n int) error {
	h := Card(nil, NeedResolved())
	return RefuseInvented(n, NeedResolved(), h.Copy)
}
