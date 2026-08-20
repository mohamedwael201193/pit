package redteam

import (
	"fmt"

	"github.com/mohamedwael201193/pit/internal/session"
)

var ExtraDenied = []string{
	"withdraw3", "updateLeverage", "sendAsset", "approveAgent",
	"usdSend", "spotSend", "vaultTransfer", "twapOrder", "modify",
}

func SessionActionMatrix() []error {
	var errs []error
	for _, a := range ExtraDenied {
		if err := session.CheckAction(a); err == nil {
			errs = append(errs, fmt.Errorf("allowlisted_kill:%s", a))
		}
	}
	for _, a := range []string{"order", "cancel"} {
		if err := session.CheckAction(a); err != nil {
			errs = append(errs, fmt.Errorf("blocked_live:%s", a))
		}
	}
	return errs
}
