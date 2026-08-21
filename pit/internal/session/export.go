package session

import "fmt"

func ExportDenied(surface string) error {
	switch surface {
	case "browser", "mcp", "vercel", "web":
		return fmt.Errorf("session_export_denied")
	default:
		return fmt.Errorf("session_export_denied")
	}
}
