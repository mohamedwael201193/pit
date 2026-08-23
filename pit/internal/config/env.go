package config

import (
	"fmt"
	"os"
	"strings"
)

func RefuseSessionEnv() error {
	if strings.ToLower(os.Getenv("PIT_PRODUCT_MODE")) != "true" {
		return nil
	}
	for _, k := range []string{"PIT_SESSION_KEY", "HL_SECRET", "PIT_MEMORY_KEY"} {
		if strings.TrimSpace(os.Getenv(k)) != "" {
			return fmt.Errorf("session_env_forbidden")
		}
	}
	return nil
}
