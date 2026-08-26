package redteam

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/session"
)

func TwoUsers(sessionWS, callerA, callerB string) error {
	if err := session.BindWorkspace(sessionWS, callerA); err != nil {
		return err
	}
	if err := session.BindWorkspace(sessionWS, callerB); err == nil {
		return fmt.Errorf("user_b_must_not_bind")
	}
	return nil
}
