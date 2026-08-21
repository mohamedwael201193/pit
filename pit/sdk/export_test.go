package sdk

import "testing"

func TestExportSessionDenied(t *testing.T) {
	c := Client{}
	if err := c.ExportSession(); err == nil {
		t.Fatal("export")
	}
}
