package session

import "testing"

func TestExportDeniedEverywhere(t *testing.T) {
	for _, s := range []string{"browser", "mcp", "desktop-log", "cli-stdout"} {
		if err := ExportDenied(s); err == nil {
			t.Fatal(s)
		}
	}
}
