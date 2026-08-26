package sdk

import "github.com/mohamedwael201193/pit/internal/session"

func BindWorkspace(sessionWS, callerWS string) error {
	return session.BindWorkspace(sessionWS, callerWS)
}
