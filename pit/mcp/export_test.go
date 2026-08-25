package mcp

import "testing"

func TestExportNever(t *testing.T) {
	if !ExportNever() {
		t.Fatal("export")
	}
}
