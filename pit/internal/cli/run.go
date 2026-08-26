package cli

import "github.com/mohamedwael201193/pit/internal/session"

func RunAuthorize(isTTY bool, typed string, iUnderstand bool, s session.Session, nowMs int64) error {
	return ConfirmAuthorizeSession(isTTY, typed, iUnderstand, session.Alive(s, nowMs))
}
