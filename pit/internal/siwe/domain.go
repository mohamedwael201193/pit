package siwe

import (
	"fmt"
	"strings"
)

func DomainOK(domain string) error {
	d := strings.ToLower(strings.TrimSpace(domain))
	if d == "" {
		return fmt.Errorf("incomplete siwe")
	}
	if strings.Contains(d, " ") {
		return fmt.Errorf("incomplete siwe")
	}
	return nil
}
