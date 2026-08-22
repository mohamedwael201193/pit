package exec

import "github.com/mohamedwael201193/pit/internal/session"

func RefuseWithdraw(action string) error {
	return session.CheckAction(action)
}
