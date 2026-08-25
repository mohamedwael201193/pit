package cli

import "github.com/mohamedwael201193/pit/internal/session"

func CancelOnly(action string) error {
	return session.CheckAction(action)
}
