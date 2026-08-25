package chain8004

import (
	"fmt"
	"strings"
)

func RefusePrivateTag(tag string) error {
	low := strings.ToLower(tag)
	for _, b := range banned {
		if strings.Contains(low, b) {
			return fmt.Errorf("private_field_forbidden")
		}
	}
	return nil
}
