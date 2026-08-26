package storage

import (
	"fmt"
	"strings"
)

func RefuseGlobalMemoryKey(v string) error {
	if strings.TrimSpace(v) != "" {
		return fmt.Errorf("global_memory_key_denied")
	}
	return nil
}
