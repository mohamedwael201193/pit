package cli

import (
	"fmt"
	"strings"
)

const ConfirmToken = "AUTHORIZE"

// ConfirmAuthorize requires an interactive TTY and the exact token AUTHORIZE.
// Piped "yes" is never enough. --i-understand still requires the token.
func ConfirmAuthorize(isTTY bool, typed string, iUnderstand bool) error {
	if !isTTY {
		return fmt.Errorf("piped_authorize_denied")
	}
	if strings.TrimSpace(typed) != ConfirmToken {
		return fmt.Errorf("need_exact_AUTHORIZE")
	}
	if !iUnderstand {
		return fmt.Errorf("need_i_understand_flag")
	}
	return nil
}

func ConfirmAuthorizeSession(isTTY bool, typed string, iUnderstand, sessionAlive bool) error {
	if !sessionAlive {
		return fmt.Errorf("session_expired")
	}
	return ConfirmAuthorize(isTTY, typed, iUnderstand)
}

func StdinIsTTY(mode uint32, isCharDevice bool) bool {
	return isCharDevice && mode != 0
}
