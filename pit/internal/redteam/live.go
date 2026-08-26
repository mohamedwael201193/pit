package redteam

import (
	"github.com/mohamedwael201193/pit/internal/cli"
	"github.com/mohamedwael201193/pit/internal/session"
)

func ExpiredAuthorize() error {
	return cli.RunAuthorize(true, "AUTHORIZE", true, session.Session{}, 1)
}
