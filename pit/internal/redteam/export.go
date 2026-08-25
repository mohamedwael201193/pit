package redteam

import "github.com/mohamedwael201193/pit/internal/session"

func SessionExportDenied() bool {
	return session.CheckAction("export_session") != nil
}
