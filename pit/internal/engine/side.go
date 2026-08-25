package engine

import (
	"fmt"
	"strings"
)

func RefuseEmptySide(side string) error {
	s := strings.ToLower(strings.TrimSpace(side))
	if s != "buy" && s != "sell" {
		return fmt.Errorf("bad_side")
	}
	return nil
}
