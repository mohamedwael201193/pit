package exec

import "github.com/mohamedwael201193/pit/internal/session"

func Gate(s session.Session, nowMs int64, policyVer, workspace string) error {
	if err := session.CheckSession(s, nowMs, policyVer, workspace); err != nil {
		return err
	}
	return nil
}
