package hl

import (
	"fmt"
	"strings"
)

func ValidCloid(c string) error {
	s := strings.ToLower(strings.TrimSpace(c))
	if !strings.HasPrefix(s, "0x") || len(s) != 34 {
		return fmt.Errorf("bad_cloid")
	}
	for _, r := range s[2:] {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return fmt.Errorf("bad_cloid")
		}
	}
	return nil
}
