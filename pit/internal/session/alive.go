package session

func Alive(s Session, nowMs int64) bool {
	return CheckSession(s, nowMs, s.PolicyVer, s.Workspace) == nil
}
