package mcp

import "github.com/mohamedwael201193/pit/internal/session"

func Bind(sessionWS, callerWS string) error {
	return session.BindWorkspace(sessionWS, callerWS)
}
