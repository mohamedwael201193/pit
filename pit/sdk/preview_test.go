package sdk

import "testing"

func TestPreviewFieldsBound(t *testing.T) {
	f := Client{}.PreviewFields()
	if len(f) < 12 {
		t.Fatalf("%v", f)
	}
	seen := map[string]bool{}
	for _, x := range f {
		seen[x] = true
	}
	for _, need := range []string{"market", "side", "sz", "cloid", "forecastId"} {
		if !seen[need] {
			t.Fatal(need)
		}
	}
}
