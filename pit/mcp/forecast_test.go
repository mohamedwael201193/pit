package mcp

import "testing"

func TestForecastNeverSizes(t *testing.T) {
	if !ForecastNeverSizes() {
		t.Fatal("size")
	}
}
