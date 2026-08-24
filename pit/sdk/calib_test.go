package sdk

import (
	"strings"
	"testing"
)

func TestHealthCardEmpty(t *testing.T) {
	h := Client{}.HealthCard(30)
	if h.Sufficient || !strings.Contains(strings.ToLower(h.Copy), "not enough") {
		t.Fatal(h)
	}
}
