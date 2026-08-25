package wallet

import (
	"fmt"
	"strings"
)

func RefuseStaleAddress(bound, current string) error {
	if strings.TrimSpace(bound) == "" || strings.TrimSpace(current) == "" {
		return fmt.Errorf("address_required")
	}
	if !strings.EqualFold(bound, current) {
		return fmt.Errorf("account_switched")
	}
	return nil
}
