package obs

import (
	"fmt"
	"strings"
)

func RefusePrivateBook(fields []string) error {
	for _, f := range fields {
		low := strings.ToLower(f)
		if low == "book" || low == "strategy" || low == "positions" {
			return fmt.Errorf("private_field_in_log")
		}
	}
	return nil
}
