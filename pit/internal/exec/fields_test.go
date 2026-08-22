package exec

import "testing"

func TestMissingBindField(t *testing.T) {
	present := map[string]bool{}
	for _, f := range PreviewBindFields {
		present[f] = true
	}
	if MissingBindField(present) != "" {
		t.Fatal("complete")
	}
	delete(present, "sz")
	if MissingBindField(present) != "sz" {
		t.Fatal("sz")
	}
}
