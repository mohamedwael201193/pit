package obs

import (
	"strings"
	"testing"
)

func TestNewRequestID(t *testing.T) {
	a := NewRequestID()
	b := NewRequestID()
	if a == b || !strings.HasPrefix(a, "req_") {
		t.Fatalf("%s %s", a, b)
	}
}
