package cli

import "github.com/mohamedwael201193/pit/internal/session"

func RefuseLeverageChange(action string) error {
	return session.CheckAction(action)
}
