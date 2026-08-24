package obs

import "testing"

func TestLatency(t *testing.T) {
	if err := Latency(12); err != nil {
		t.Fatal(err)
	}
	if err := Latency(-1); err == nil {
		t.Fatal("neg")
	}
}
