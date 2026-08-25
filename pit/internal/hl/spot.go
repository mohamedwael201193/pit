package hl

func SpotCountsAsFunded(v AccountView) bool {
	return v.State == FundedSpot || v.State == FundedPerp
}
