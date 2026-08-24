package mcp

import "testing"

func TestTransferNever(t *testing.T) {
	if !TransferNever() {
		t.Fatal("xfer")
	}
}
