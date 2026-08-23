package siwe

import (
	"fmt"
	"strings"
)

func OriginMatch(expected, got string) error {
	if err := DomainOK(expected); err != nil {
		return err
	}
	if err := DomainOK(got); err != nil {
		return err
	}
	if !strings.EqualFold(strings.TrimSpace(expected), strings.TrimSpace(got)) {
		return fmt.Errorf("origin_mismatch")
	}
	return nil
}
