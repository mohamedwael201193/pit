package memory

import (
	"fmt"
	"os"
)

func RefuseGlobalKey() error {
	if os.Getenv("PIT_PRODUCT_MODE") == "true" && os.Getenv("PIT_MEMORY_KEY") != "" {
		return fmt.Errorf("global_memory_key_forbidden")
	}
	return nil
}
