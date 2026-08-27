package exec

import "testing"

func TestRefusePostUntilLinked(t *testing.T) {
	if err := RefusePostUntilLinked(false); err == nil {
		t.Fatal("unlinked")
	}
	if err := RefusePostUntilLinked(true); err != nil {
		t.Fatal(err)
	}
}
