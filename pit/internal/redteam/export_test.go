package redteam

import "testing"

func TestSessionExportDenied(t *testing.T) {
	if !SessionExportDenied() {
		t.Fatal("export")
	}
}
