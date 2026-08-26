package policy

import "fmt"

func RequireSession(c Context) error {
	if !c.SessionAlive {
		return fmt.Errorf("session_expired")
	}
	return nil
}
