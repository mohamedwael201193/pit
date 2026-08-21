package sdk

import (
	"github.com/mohamedwael201193/pit/internal/session"
)

func (c Client) ExportSession() error {
	return session.ExportDenied("browser")
}
