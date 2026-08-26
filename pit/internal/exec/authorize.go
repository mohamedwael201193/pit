package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/session"
)

func RequireLive(s session.Session, nowMs int64) error {
	if !session.Alive(s, nowMs) {
		return fmt.Errorf("session_expired")
	}
	return nil
}
