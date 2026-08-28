package compute

import (
	"strings"
	"testing"
)

func TestThreeEnvelopesSameProviderHonest(t *testing.T) {
	got, err := Committee([]byte(`{"eth":1}`), []byte(`{"clip":10}`))
	if err != nil || len(got) != 3 {
		t.Fatal(err, got)
	}
	if IndependenceNote() != EnvelopeOnly {
		t.Fatal("must not claim independent providers")
	}
	if string(got[Researcher]) == string(got[Challenger]) {
		t.Fatal("envelopes must differ by role")
	}
	if !strings.Contains(string(got[Researcher]), `"proposed_side"`) {
		t.Fatal("researcher must ask for JSON side")
	}
}
