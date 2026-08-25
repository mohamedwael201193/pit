package cli

import (
	"github.com/mohamedwael201193/pit/internal/exec"
	"github.com/mohamedwael201193/pit/internal/session"
)

func RefuseKilledSession(kill bool) error {
	return exec.RefuseKill(session.Session{Kill: kill})
}
