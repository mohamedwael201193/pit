package policy

func InCooldown(lastFillUnix, nowUnix, seconds int64) bool {
	if seconds <= 0 || lastFillUnix <= 0 || nowUnix <= 0 {
		return false
	}
	return nowUnix-lastFillUnix < seconds
}
