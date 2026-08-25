package phase

import "testing"

func TestWebMaySign(t *testing.T) {
	if WebMaySign("web") || WebMaySign("browser") {
		t.Fatal("web")
	}
	if !WebMaySign("desktop") || !WebMaySign("cli") {
		t.Fatal("host")
	}
}
