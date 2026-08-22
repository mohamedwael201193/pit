package keyring

import (
	"fmt"
	"strings"
)

func RefuseWeb(surface string) error {
	s := strings.ToLower(strings.TrimSpace(surface))
	switch s {
	case "web", "browser", "vercel", "mcp":
		return fmt.Errorf("keyring_not_in_browser")
	default:
		return nil
	}
}

func Redact(secret []byte) string {
	if len(secret) == 0 {
		return ""
	}
	return "[redacted]"
}
