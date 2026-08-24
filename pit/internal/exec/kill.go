package exec

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/session"
)

func RefuseKill(s session.Session) error {
	if s.Kill {
		return fmt.Errorf("kill_switch")
	}
	return nil
}
